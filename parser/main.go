package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// Inspired by https://github.com/golang/go/blob/master/src/text/template/parse/lex.go
// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos Pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

// https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#string_literals
// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF                   // end

	// letters and digits
	itemLetter        // "A" … "Z" | "a" … "z"
	itemCapitalLetter //  "A" … "Z"
	itemDecimalDigit  // "0" … "9"
	itemOctalDigit    // "0" … "7"
	itemHexDigit      // "0" … "9" | "A" … "F" | "a" … "f"

	// identifiers
	itemIdent       // letter { letter | unicodeDigit | "_" }
	itemFullIdent   // ident { "." ident }
	itemMessageName // ident
	itemEnumName    // ident
	itemFieldName   // ident
	itemOneofName   // ident
	itemMapName     // ident
	itemServiceName // ident
	itemRPCName     // ident
	itemMessageType // [ "." ] { ident "." } messageName
	itemEnumType    // [ "." ] { ident "." } enumName

	// integer literals
	itemIntLit     // decimalLit | octalLit | hexLit
	itemDecimalLit // ( "1" … "9" ) { decimalDigit }
	itemOctalLit   // "0" { octalDigit }
	itemHexLit     // "0" ( "x" | "X" ) hexDigit { hexDigit }

	// floating-point literals
	itemFloatLit // decimals "." [ decimals ] [ exponent ] | decimals exponent | "."decimals [ exponent ]
	itemDecimals // decimalDigit { decimalDigit }
	itemExponent // ( "e" | "E" ) [ "+" | "-" ] decimals

	// boolean
	itemBoolLit // "true" | "false"

	// string literals
	itemStrLit     // ( "'" { charValue } "'" ) |  ( '"' { charValue } '"' )
	itemCharValue  // hexEscape | octEscape | charEscape | /[^\0\n\\]/
	itemHexEscape  // '\' ( "x" | "X" ) hexDigit hexDigit
	itemOctEscape  // '\' octalDigit octalDigit octalDigit
	itemCharEscape // '\' ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | '\' | "'" | '"' )
	itemQuote      // "'" | '"'

	// empty
	itemEmptyStatement // ";"
	itemConstant       // fullIdent | ( [ "-" | "+" ] intLit ) | ( [ "-" | "+" ] floatLit ) | strLit | boolLit
	itemSyntax         // "syntax" "=" quote "proto3" quote ";"
	itemImport         // "import" [ "weak" | "public" ] strLit ";"
	itemPackage        // "package" fullIdent ";"

	// option
	itemOption     // "option" optionName  "=" constant ";"
	itemOptionName // ( ident | "(" fullIdent ")" ) { "." ident }

	// fields
	itemFieldType   // "double" | "float" | "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string" | "bytes" | messageType | enumType
	itemFieldNumber // intLit;

	itemField        // [ "repeated" ] type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
	itemFieldOptions // fieldOption { ","  fieldOption }
	itmeFieldOption  // optionName "=" constant

	itemOneOf      // "oneof" oneofName "{" { oneofField | emptyStatement } "}"
	itemOneOfField // type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"

	itemMapField // "map" "<" keyType "," type ">" mapName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
	itemKeyType  // "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string"

	itemReserved   // "reserved" ( ranges | fieldNames ) ";"
	itemFieldNames // fieldName { "," fieldName }

	// top level
	itemEnum            // "enum" enumName enumBody
	itemEnumBody        // "{" { option | enumField | emptyStatement } "}"
	itemEnumField       // ident "=" intLit [ "[" enumValueOption { ","  enumValueOption } "]" ]";"
	itemEnumValueOption // optionName "=" constant

	itemMessage     // "message" messageName messageBody
	itemMessageBody // "{" { field | enum | message | option | oneof | mapField | reserved | emptyStatement } "}"

	itemService // "service" serviceName "{" { option | rpc | stream | emptyStatement } "}"
	itemRPC     // "rpc" rpcName "(" [ "stream" ] messageType ")" "returns" "(" [ "stream" ] messageType // ")" (( "{" {option | emptyStatement } "}" ) | ";")

	itemProto       // syntax { import | package | option | topLevelDef | emptyStatement }
	itemTopLevelDef // message | enum | service
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

const eof = -1

const (
	tokSyntax = "syntax"
)

// lexer holds the state of the scanner.
type lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	leftDelim  string    // start of action
	rightDelim string    // end of action
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	lastPos    Pos       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexSyntax; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.pos, l.input[l.start:l.pos]}
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// Consume spaces
func (l *lexer) trim() {
	for isSpace(l.peek()) {
		l.next()
	}
}

func lexSyntax(l *lexer) stateFn {
	for _, tok := range []string{"syntax", "=", "\"proto3\";"} {
		l.trim()
		if strings.HasPrefix(l.input[l.pos:], tok) {
			l.pos += Pos(len(tok))
			continue
		}
		return l.errorf("proto file must start with 'syntax = \"proto3\";")
	}
	l.emit(itemSyntax)
	return lexEnd
}

func lexEnd(l *lexer) stateFn {
	l.emit(itemEOF)
	return nil
}
