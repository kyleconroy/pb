package ast

type syntax int

const (
	Proto2 syntax = iota
	Proto3
)

type File struct {
	Syntax syntax
}
