package parser

import (
	"errors"
	"io"
	"io/ioutil"
	"log"

	"github.com/kyleconroy/pb/ast"
	"github.com/kyleconroy/pb/token"
)

type Mode int

func ParseFile(fset *token.FileSet, filename string, src io.Reader, mode Mode) (*ast.File, error) {
	payload, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	l := lex(filename, string(payload))
	defer l.drain()
	for i := range l.items {
		log.Println(i)
		if i.typ == itemError {
			return nil, errors.New(i.val)
		}
	}
	return &ast.File{}, nil
}
