package parser

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/kyleconroy/pb/ast"
)

type Mode int

func ParseFile(src io.Reader, mode Mode) (*ast.File, error) {
	payload, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	t := tree{lex("", string(payload)), &ast.File{}}
	return t.parse()
}

type tree struct {
	l *lexer
	f *ast.File
}

func (t *tree) parse() (*ast.File, error) {
	defer t.l.drain()

	if err := t.parseSyntax(); err != nil {
		return t.f, err
	}

	//for {
	//	token := t.nextNonComment()
	//	if token.typ == itemError {
	//		return t.f, errors.New(token.val)
	//	}
	//}

	return t.f, nil
}

func (t *tree) parseSyntax() error {
	for _, f := range []func(item) bool{
		func(i item) bool { return i.typ == itemSyntax },
		func(i item) bool { return i.typ == itemEq },
		func(i item) bool { return i.typ == itemStrLit && i.val == "\"proto3\"" },
		func(i item) bool { return i.typ == itemSemiColon },
	} {
		item := t.nextNonComment()
		if !f(item) {
			return errors.New("proto files must start with `syntax =\"proto3\";`")
		}
	}
	t.f.Syntax = ast.Proto3
	return nil
}

// nextNonSpace returns the next non-space token.
func (t *tree) nextNonComment() (token item) {
	for {
		token = t.l.nextItem()
		if token.typ != itemComment {
			break
		}
	}
	return token
}
