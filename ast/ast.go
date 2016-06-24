package ast

import (
	"github.com/kyleconroy/pb/token"
)

type syntax int

const (
	Proto2 syntax = iota
	Proto3
)

type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

type BasicLit struct {
	Kind  token.Token
	Value string
}

func (bs *BasicLit) Pos() token.Pos {
	return token.Pos(0)
}

func (bs *BasicLit) End() token.Pos {
	return token.Pos(0)
}

type BlockStmt struct {
	Opening token.Pos // position of opening parenthesis/brace, if any
	List    []Node
	Closing token.Pos // position of closing parenthesis/brace, if any
}

func (s *BlockStmt) Pos() token.Pos {
	return s.Opening
}

func (s *BlockStmt) End() token.Pos {
	return s.Closing
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
	Name  *Ident
	Value string
}

func (e *EnumField) Pos() token.Pos {
	return token.Pos(0)
}

func (e *EnumField) End() token.Pos {
	return token.Pos(0)
}

type Expr struct {
}

func (e *Expr) Pos() token.Pos {
	return token.Pos(0)
}

func (e *Expr) End() token.Pos {
	return token.Pos(0)
}

type File struct {
	Syntax syntax
	Nodes  []Node
}

func (f *File) Pos() token.Pos {
	return token.Pos(0)
}

func (f *File) End() token.Pos {
	return token.Pos(0)
}

type Import struct {
	Modifiers []*Ident
	Path      *BasicLit
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

type MapType struct {
	Map   token.Pos // position of "map" keyword
	Key   *Ident
	Value *Ident
}

func (m *MapType) Pos() token.Pos {
	return token.Pos(0)
}

func (m *MapType) End() token.Pos {
	return token.Pos(0)
}

type Message struct {
	Name *Ident
	Body []Node
}

func (m *Message) Pos() token.Pos {
	return token.Pos(0)
}

func (m *Message) End() token.Pos {
	return token.Pos(0)
}

type MessageField struct {
	Name     *Ident
	Number   *BasicLit
	Type     Node
	Repeated *Ident
}

func (m *MessageField) Pos() token.Pos {
	return token.Pos(0)
}

func (m *MessageField) End() token.Pos {
	return token.Pos(0)
}

type OneOf struct {
	Name  *Ident
	Body  []Node
	OneOf token.Pos // position of "oneof" keyword
}

func (oo *OneOf) Pos() token.Pos {
	return token.Pos(0)
}

func (oo *OneOf) End() token.Pos {
	return token.Pos(0)
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

type Package struct {
}

func (p *Package) Pos() token.Pos {
	return token.Pos(0)
}

func (p *Package) End() token.Pos {
	return token.Pos(0)
}

type RPC struct {
	RPC     token.Pos
	Name    Ident
	InType  Ident
	OutType Ident
}

func (r *RPC) Pos() token.Pos {
	return token.Pos(0)
}

func (r *RPC) End() token.Pos {
	return token.Pos(0)
}

type Service struct {
	Service token.Pos
	Name    Ident
	Body    *BlockStmt
}

func (s *Service) Pos() token.Pos {
	return token.Pos(0)
}

func (s *Service) End() token.Pos {
	return token.Pos(0)
}
