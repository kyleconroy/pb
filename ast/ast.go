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
	Body int
}

func (m *Message) Pos() token.Pos {
	return token.Pos(0)
}

func (m *Message) End() token.Pos {
	return token.Pos(0)
}

type Enum struct {
	Name Ident
	Body int
}

func (e *Enum) Pos() token.Pos {
	return token.Pos(0)
}

func (e *Enum) End() token.Pos {
	return token.Pos(0)
}

type Service struct {
	Name Ident
	Body int
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
