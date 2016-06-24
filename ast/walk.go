package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {

	case *BasicLit:

	case *BlockStmt:
		for _, m := range n.List {
			Walk(v, m)
		}

	case *EmptyStmt, *Expr:

	case *Enum:
		for _, m := range n.Body {
			Walk(v, m)
		}

	case *EnumField:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	case *File:
		for _, m := range n.Nodes {
			Walk(v, m)
		}

	case *Ident:

	case *Import:
		for _, m := range n.Modifiers {
			Walk(v, m)
		}
		if n.Path != nil {
			Walk(v, n.Path)
		}

	case *MapType:
		if n.Key != nil {
			Walk(v, n.Key)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *Message:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, m := range n.Body {
			Walk(v, m)
		}

	case *MessageField:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Number != nil {
			Walk(v, n.Number)
		}
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Repeated != nil {
			Walk(v, n.Repeated)
		}

	case *OneOf:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, m := range n.Body {
			Walk(v, m)
		}

	case *Option:
		for _, m := range n.Names {
			Walk(v, m)
		}
		if n.Constant != nil {
			Walk(v, n.Constant)
		}

	case *Package:
		// TODO: Add fields to Package

	case *RPC:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.InType != nil {
			Walk(v, n.InType)
		}
		if n.OutType != nil {
			Walk(v, n.OutType)
		}

	case *Service:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}
