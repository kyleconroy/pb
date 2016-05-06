package parser

import (
	"fmt"
	"unicode"
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

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#string_literals
// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF                   // end

	itemDot          // .
	itemEq           // =
	itemSemiColon    // ;
	itemLeftBrace    // {
	itemRightBrace   // }
	itemLeftParen    // (
	itemRightParen   // )
	itemLeftBracket  // [
	itemRightBracket // ]
	itemLeftMap      // <
	itemRightMap     // >
	itemIdent        // letter { letter | unicodeDigit | "_" }
	itemFullIdent    // ident { "." ident }
	itemStrLit       // ( "'" { charValue } "'" ) |  ( '"' { charValue } '"' )
	itemComment      // // comment

	// keywords
	itemKeyword
	itemSyntax       // syntax
	itemMessage      // message
	itemEnum         // enum
	itemImport       // import
	itemImportWeak   // weak
	itemImportPublic // public
	itemPackage      // package
	itemOption       // option
	itemMap          // map
	itemRepeated     // repreated
	itemReturns      // returns
	itemRPC          // rpc
	itemService      // service
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

const eof = -1

const (
	tokSyntax = "syntax"
)

// lexer holds the state of the scanner.
type lexer struct {
	name         string    // the name of the input; used only for error reports
	input        string    // the string being scanned
	leftDelim    string    // start of action
	rightDelim   string    // end of action
	state        stateFn   // the next lexing function to enter
	pos          Pos       // current position in the input
	start        Pos       // start position of this item
	width        Pos       // width of last rune read from input
	lastPos      Pos       // position of most recent item returned by nextItem
	items        chan item // channel of scanned items
	braceDepth   int       // nesting depth of { }
	bracketDepth int       // nesting depth of [ ]
	parenDepth   int       // nesting depth of ( )
	mapDepth     int       // nesting depth of < >
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexSchema; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for _ = range l.items {
	}
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
	return r == ' ' || r == '\t' || r == '\n'
}

// Consume spaces
func (l *lexer) trim() {
	for isSpace(l.peek()) {
		l.next()
	}
	l.start = l.pos
}

func lexSchema(l *lexer) stateFn {
	// Ignore whitespace, it doesn't matter
	l.trim()
	switch r := l.next(); {
	case r == '.':
		l.emit(itemDot)
	case r == '<':
		l.emit(itemLeftMap)
		l.mapDepth++
	case r == '>':
		l.emit(itemRightMap)
		l.mapDepth--
	case r == '[':
		l.emit(itemLeftBracket)
		l.bracketDepth++
	case r == ']':
		l.emit(itemRightBracket)
		l.bracketDepth--
	case r == '(':
		l.emit(itemLeftParen)
		l.parenDepth++
	case r == ')':
		l.emit(itemRightParen)
		l.parenDepth--
	case r == '{':
		l.emit(itemLeftBrace)
		l.braceDepth++
	case r == '}':
		l.emit(itemRightBrace)
		l.braceDepth--
		if l.braceDepth < 0 {
			return l.errorf("unexpected right brace %#U", r)
		}
	case r == ',':
		if l.mapDepth <= 0 {
			return l.errorf("unexpected comman outside of map definition %#U", r)
		}
		l.emit(itemRightBrace)
	case r == '=':
		l.emit(itemEq)
	case r == '"':
		return lexQuote
	case r == ';':
		l.emit(itemSemiColon)
	case r == '/':
		return lexComment
	case isAlphaNumeric(r):
		l.backup()
		return lexIdent
	case r == eof:
		return lexEnd
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexSchema
}

var key = map[string]itemType{
	"syntax":   itemSyntax,
	"import":   itemImport,
	"weak":     itemImportWeak,
	"public":   itemImportPublic,
	"message":  itemMessage,
	"enum":     itemEnum,
	"option":   itemOption,
	"map":      itemMap,
	"rpc":      itemRPC,
	"returns":  itemReturns,
	"service":  itemService,
	"repeated": itemRepeated,
}

func lexComment(l *lexer) stateFn {
	if l.next() != '/' {
		return l.errorf("comments must start with two backslashes")
	}
	for {
		r := l.next()
		if r == '\n' || r == '\r' || r == eof {
			l.emit(itemComment)
			return lexSchema
		}
	}
}

// lexIdentifier scans an alphanumeric.
func lexIdent(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
			case key[word] != itemError:
				l.emit(key[word])
				return lexSchema
			default:
				l.emit(itemIdent)
				return lexSchema
			}
		}
	}
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) stateFn {
	for {
		switch l.next() {
		//case '\\':
		//	if r := l.next(); r != eof && r != '\n' {
		//		break
		//	}
		//	fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			l.emit(itemStrLit)
			return lexSchema
		}
	}
}

func lexEnd(l *lexer) stateFn {
	l.emit(itemEOF)
	return nil
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
