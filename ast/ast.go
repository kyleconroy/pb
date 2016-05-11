package ast

import (
	"github.com/kyleconroy/pb/token"
)

type syntax int

const (
	Proto2 syntax = iota
	Proto3
)

type File struct {
	Syntax syntax
	Nodes  []Node
}

type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

type Import struct {
	Modifiers []Ident
	Path      BasicLit
}

type BasicLit struct {
	Kind  token.Token
	Value string
}

func (i *Import) Pos() token.Pos {
	return token.Pos(0)
}

func (i *Import) End() token.Pos {
	return token.Pos(0)
}

type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name
}

func (i *Ident) Pos() token.Pos {
	return token.Pos(0)
}

func (i *Ident) End() token.Pos {
	return token.Pos(0)
}

type Expr struct {
}

type Option struct {
	Names    []Ident
	Constant BasicLit
}

func (o *Option) Pos() token.Pos {
	return token.Pos(0)
}

func (o *Option) End() token.Pos {
	return token.Pos(0)
}

type Message struct {
	Name Ident
	Body []Node
}

func (m *Message) Pos() token.Pos {
	return token.Pos(0)
}

func (m *Message) End() token.Pos {
	return token.Pos(0)
}

type MapType struct {
	Map   token.Pos // position of "map" keyword
	Key   Ident
	Value Ident
}

func (m *MapType) Pos() token.Pos {
	return token.Pos(0)
}

func (m *MapType) End() token.Pos {
	return token.Pos(0)
}

type MessageField struct {
	Name     Ident
	Number   BasicLit
	Type     Node
	Repeated *Ident
}

func (m *MessageField) Pos() token.Pos {
	return token.Pos(0)
}

func (m *MessageField) End() token.Pos {
	return token.Pos(0)
}

type EmptyStmt struct {
	Semicolon token.Pos // position of following ";"
}

func (e *EmptyStmt) Pos() token.Pos {
	return token.Pos(0)
}

func (e *EmptyStmt) End() token.Pos {
	return token.Pos(0)
}

type Enum struct {
	Name Ident
	Body []Node
}

func (e *Enum) Pos() token.Pos {
	return token.Pos(0)
}

func (e *Enum) End() token.Pos {
	return token.Pos(0)
}

type EnumField struct {
	Name  Ident
	Value string
}

func (e *EnumField) Pos() token.Pos {
	return token.Pos(0)
}

func (e *EnumField) End() token.Pos {
	return token.Pos(0)
}

type Service struct {
	Name Ident
	Body []Node
}

func (s *Service) Pos() token.Pos {
	return token.Pos(0)
}

func (s *Service) End() token.Pos {
	return token.Pos(0)
}

type Package struct {
}

func (p *Package) Pos() token.Pos {
	return token.Pos(0)
}

func (p *Package) End() token.Pos {
	return token.Pos(0)
}
