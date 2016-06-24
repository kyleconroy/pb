package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kyleconroy/pb/ast"
	"github.com/kyleconroy/pb/token"
)

type Mode int

func ParseFile(src io.Reader, mode Mode) (*ast.File, error) {
	payload, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	t := tree{lex("", string(payload)), &ast.File{Nodes: []ast.Node{}}}
	return t.parse()
}

type tree struct {
	l *lexer
	f *ast.File
}

func (t *tree) expect(typs ...itemType) ([]item, error) {
	items := make([]item, len(typs))
	for i, typ := range typs {
		tok := t.nextNonComment()
		if tok.typ == typ {
			items[i] = tok
		} else {
			return items, fmt.Errorf("Incorrect token: %s", tok)
		}
	}
	return items, nil
}

func (t *tree) parse() (*ast.File, error) {
	defer t.l.drain()

	if err := t.parseSyntax(); err != nil {
		return t.f, err
	}

	for {
		switch token := t.nextNonComment(); {
		case token.typ == itemImport:
			if err := t.parseImport(); err != nil {
				return t.f, err
			}
		case token.typ == itemPackage:
			if err := t.parsePackage(); err != nil {
				return t.f, err
			}
		case token.typ == itemOption:
			if err := t.parseOption(); err != nil {
				return t.f, err
			}
		case token.typ == itemMessage:
			node, err := t.parseMessage()
			if err != nil {
				return t.f, err
			}
			t.f.Nodes = append(t.f.Nodes, node)
		case token.typ == itemService:
			node, err := t.parseService(token)
			if err != nil {
				return t.f, err
			}
			t.f.Nodes = append(t.f.Nodes, node)
		case token.typ == itemEnum:
			node, err := t.parseEnum()
			if err != nil {
				return t.f, err
			}
			t.f.Nodes = append(t.f.Nodes, node)
		case token.typ == itemError:
			return t.f, errors.New(token.val)
		case token.typ == itemEOF:
			return t.f, nil
		default:
			return t.f, fmt.Errorf("Incorrect token: %s", token)
		}
	}

	return t.f, nil
}

func (t *tree) parseSyntax() error {
	toks, err := t.expect(itemSyntax, itemEq, itemStrLit, itemSemiColon)
	if err != nil {
		return err
	}
	if toks[2].val != "\"proto3\"" {
		return fmt.Errorf("Looking for proto3")
	}
	t.f.Syntax = ast.Proto3
	return nil
}

func (t *tree) parsePackage() error {
	for {
		item := t.nextNonComment()
		if item.typ == itemSemiColon {
			break
		}
	}
	t.f.Nodes = append(t.f.Nodes, &ast.Package{})
	return nil
}

func (t *tree) parseImport() error {
	idents := []*ast.Ident{}
	seen := map[itemType]struct{}{}
	for {
		switch tok := t.nextNonComment(); {
		case tok.typ == itemImportPublic || tok.typ == itemImportWeak:
			if _, ok := seen[tok.typ]; ok {
				return fmt.Errorf("Multiple %s modifiers found", tok.val)
			}
			seen[tok.typ] = struct{}{}
			idents = append(idents, &ast.Ident{Name: tok.val})
		case tok.typ == itemStrLit:
			if end := t.nextNonComment(); end.typ != itemSemiColon {
				return fmt.Errorf("Incorrect token: %s", end)
			}
			t.f.Nodes = append(t.f.Nodes, &ast.Import{
				Modifiers: idents,
				Path:      &ast.BasicLit{Value: tok.val, Kind: token.STRING},
			})
			return nil
		default:
			return fmt.Errorf("Incorrect token: %s", tok)
		}
	}
	return nil
}

func (t *tree) parseOption() error {
	var ident item

	tok := t.nextNonComment()
	// TODO We need to handle full idents
	if tok.typ != itemIdent {
		return fmt.Errorf("expected ident, found %s", tok)
	}
	ident = tok

	tok = t.nextNonComment()
	// TODO We need to handle full idents
	if tok.typ != itemEq {
		return fmt.Errorf("expected =, found %s", tok)
	}

	tok = t.nextNonComment()
	var con ast.BasicLit
	// TODO We need to handle all constant types
	switch tok.typ {
	case itemStrLit:
		con = ast.BasicLit{Value: tok.val, Kind: token.STRING}
	case itemBoolLit:
		con = ast.BasicLit{Value: tok.val, Kind: token.BOOL}
	default:
		return fmt.Errorf("expected string literal, found %s", tok)
	}

	if end := t.nextNonComment(); end.typ != itemSemiColon {
		return fmt.Errorf("Incorrect token: %s", end)
	}

	t.f.Nodes = append(t.f.Nodes, &ast.Option{
		Names: []*ast.Ident{
			{Name: ident.val},
		},
		Constant: &con,
	})
	return nil
}

func (t *tree) parseMessage() (ast.Node, error) {
	name := t.nextNonComment()
	if name.typ != itemIdent {
		return nil, fmt.Errorf("expected ident, found %s", name)
	}
	msg := ast.Message{
		Name: &ast.Ident{Name: name.val},
		Body: []ast.Node{},
	}

	lBrace := t.nextNonComment()
	if lBrace.typ != itemLeftBrace {
		return nil, fmt.Errorf("expected {, found %s", lBrace)
	}

	for {
		switch tok := t.nextNonComment(); {
		case tok.typ == itemSemiColon:
			msg.Body = append(msg.Body, &ast.EmptyStmt{Semicolon: token.Pos(0)})
		case tok.typ == itemOneOf:
			nmsg, err := t.parseOneOf()
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, nmsg)
		case tok.typ == itemMessage:
			nmsg, err := t.parseMessage()
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, nmsg)
		case tok.typ == itemEnum:
			nenum, err := t.parseEnum()
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, nenum)
		case tok.typ == itemOption:
			// Should be constant here, not bool
			toks, err := t.expect(itemIdent, itemEq, itemBoolLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.Option{
				Names:    []*ast.Ident{{Name: toks[0].val}},
				Constant: &ast.BasicLit{Kind: token.BOOL, Value: toks[2].val},
			})
		case tok.typ == itemRepeated:
			toks, err := t.expect(itemIdent, itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Repeated: &ast.Ident{Name: tok.val},
				Type:     &ast.Ident{Name: toks[0].val},
				Name:     &ast.Ident{Name: toks[1].val},
				Number:   &ast.BasicLit{Kind: token.INT, Value: toks[3].val},
			})
		case tok.typ == itemIdent:
			toks, err := t.expect(itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Type:   &ast.Ident{Name: tok.val},
				Name:   &ast.Ident{Name: toks[0].val},
				Number: &ast.BasicLit{Kind: token.INT, Value: toks[2].val},
			})
		case tok.typ == itemMap:
			mapt, err := t.expect(itemLeftMap, itemIdent, itemComma, itemIdent, itemRightMap)
			if err != nil {
				return nil, err
			}
			toks, err := t.expect(itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Type: &ast.MapType{
					Key:   &ast.Ident{Name: mapt[1].val},
					Value: &ast.Ident{Name: mapt[3].val},
				},
				Name:   &ast.Ident{Name: toks[0].val},
				Number: &ast.BasicLit{Kind: token.INT, Value: toks[2].val},
			})
		case tok.typ == itemRightBrace:
			return &msg, nil
		default:
			return nil, fmt.Errorf("unexpected token in message: %s", tok)
		}
	}
}

func (t *tree) parseOneOf() (ast.Node, error) {
	name := t.nextNonComment()
	if name.typ != itemIdent {
		return nil, fmt.Errorf("expected ident, found %s", name)
	}
	msg := ast.OneOf{
		Name: &ast.Ident{Name: name.val},
		Body: []ast.Node{},
	}

	lBrace := t.nextNonComment()
	if lBrace.typ != itemLeftBrace {
		return nil, fmt.Errorf("expected {, found %s", lBrace)
	}

	for {
		switch tok := t.nextNonComment(); {
		case tok.typ == itemSemiColon:
			msg.Body = append(msg.Body, &ast.EmptyStmt{Semicolon: token.Pos(0)})
		case tok.typ == itemMessage:
			nmsg, err := t.parseMessage()
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, nmsg)
		case tok.typ == itemOption:
			// Should be constant here, not bool
			toks, err := t.expect(itemIdent, itemEq, itemBoolLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.Option{
				Names:    []*ast.Ident{{Name: toks[0].val}},
				Constant: &ast.BasicLit{Kind: token.BOOL, Value: toks[2].val},
			})
		case tok.typ == itemRepeated:
			toks, err := t.expect(itemIdent, itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Repeated: &ast.Ident{Name: tok.val},
				Type:     &ast.Ident{Name: toks[0].val},
				Name:     &ast.Ident{Name: toks[1].val},
				Number:   &ast.BasicLit{Kind: token.INT, Value: toks[3].val},
			})
		case tok.typ == itemIdent:
			toks, err := t.expect(itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Type:   &ast.Ident{Name: tok.val},
				Name:   &ast.Ident{Name: toks[0].val},
				Number: &ast.BasicLit{Kind: token.INT, Value: toks[2].val},
			})
		case tok.typ == itemMap:
			mapt, err := t.expect(itemLeftMap, itemIdent, itemComma, itemIdent, itemRightMap)
			if err != nil {
				return nil, err
			}
			toks, err := t.expect(itemIdent, itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.MessageField{
				Type: &ast.MapType{
					Key:   &ast.Ident{Name: mapt[1].val},
					Value: &ast.Ident{Name: mapt[3].val},
				},
				Name:   &ast.Ident{Name: toks[0].val},
				Number: &ast.BasicLit{Kind: token.INT, Value: toks[2].val},
			})
		case tok.typ == itemRightBrace:
			return &msg, nil
		default:
			return nil, fmt.Errorf("unexpected token in message: %s", tok)
		}
	}
}

func (t *tree) parseEnum() (ast.Node, error) {
	name := t.nextNonComment()
	if name.typ != itemIdent {
		return nil, fmt.Errorf("expected ident, found %s", name)
	}
	msg := ast.Enum{
		Name: ast.Ident{Name: name.val},
		Body: []ast.Node{},
	}

	lBrace := t.nextNonComment()
	if lBrace.typ != itemLeftBrace {
		return nil, fmt.Errorf("expected {, found %s", lBrace)
	}

	for {
		switch tok := t.nextNonComment(); {
		case tok.typ == itemSemiColon:
			msg.Body = append(msg.Body, &ast.EmptyStmt{Semicolon: token.Pos(0)})
		case tok.typ == itemOption:
			// Should be constant here, not bool
			toks, err := t.expect(itemIdent, itemEq, itemBoolLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.Option{
				Names:    []*ast.Ident{{Name: toks[0].val}},
				Constant: &ast.BasicLit{Kind: token.BOOL, Value: toks[2].val},
			})
		case tok.typ == itemIdent:
			toks, err := t.expect(itemEq, itemIntLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			msg.Body = append(msg.Body, &ast.EnumField{
				Name:  &ast.Ident{Name: tok.val},
				Value: toks[1].val,
			})
		case tok.typ == itemRightBrace:
			return &msg, nil
		default:
			return nil, fmt.Errorf("unexpected token in enum: %s", tok)
		}
	}
}

func (t *tree) parseService(in item) (ast.Node, error) {
	name := t.nextNonComment()
	if name.typ != itemIdent {
		return nil, fmt.Errorf("expected ident, found %s", name)
	}

	srv := ast.Service{
		Service: token.Pos(in.pos),
		Name:    &ast.Ident{NamePos: token.Pos(name.pos), Name: name.val},
	}

	lBrace := t.nextNonComment()
	if lBrace.typ != itemLeftBrace {
		return nil, fmt.Errorf("expected {, found %s", lBrace)
	}

	blk := ast.BlockStmt{
		Opening: token.Pos(lBrace.pos),
		List:    []ast.Node{},
	}

	for {
		switch tok := t.nextNonComment(); {
		case tok.typ == itemSemiColon:
			blk.List = append(blk.List, &ast.EmptyStmt{Semicolon: token.Pos(0)})
		case tok.typ == itemOption:
			// Should be constant here, not bool
			toks, err := t.expect(itemIdent, itemEq, itemBoolLit, itemSemiColon)
			if err != nil {
				return nil, err
			}
			blk.List = append(blk.List, &ast.Option{
				Names:    []*ast.Ident{{Name: toks[0].val}},
				Constant: &ast.BasicLit{Kind: token.BOOL, Value: toks[2].val},
			})
		case tok.typ == itemRPC:
			toks, err := t.expect(itemIdent, itemLeftParen, itemIdent, itemRightParen,
				itemReturns, itemLeftParen, itemIdent, itemRightParen,
				itemLeftBrace, itemRightBrace)
			if err != nil {
				return nil, err
			}
			blk.List = append(blk.List, &ast.RPC{
				RPC:     token.Pos(tok.pos),
				Name:    &ast.Ident{Name: toks[0].val, NamePos: token.Pos(toks[0].pos)},
				InType:  &ast.Ident{Name: toks[2].val, NamePos: token.Pos(toks[2].pos)},
				OutType: &ast.Ident{Name: toks[6].val, NamePos: token.Pos(toks[6].pos)},
			})
		case tok.typ == itemRightBrace:
			blk.Closing = token.Pos(tok.pos)
			srv.Body = &blk
			return &srv, nil
		default:
			return nil, fmt.Errorf("unexpected token in enum: %s", tok)
		}
	}
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
