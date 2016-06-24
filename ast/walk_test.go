package ast

import (
	"fmt"
	"testing"
)

type walker struct {
	trip []Node
}

func (w *walker) Visit(node Node) Visitor {
	w.trip = append(w.trip, node)
	return w
}

func walk(t *testing.T, n Node, expected []Node) {
	w := &walker{}
	Walk(w, n)

	for i, n := range expected {
		if len(w.trip)-1 < i {
			t.Fatalf("Expected trip length of %d, got %d", len(expected), len(w.trip))
		}
		if n != w.trip[i] {
			t.Errorf("Expected item at index %d to be %#v, got %#v", i, n, w.trip[i])
		}
	}
	if len(expected) != len(w.trip) {
		t.Fatalf("Expected trip length of %d, got %d", len(expected), len(w.trip))
	}
}

func TestWalkLeaf(t *testing.T) {
	for _, n := range []Node{
		&BasicLit{},
		&Ident{},
		&EmptyStmt{},
		&Expr{},
	} {
		t.Run(fmt.Sprintf("%T", n), func(t *testing.T) {
			walk(t, n, []Node{n, nil})
		})
	}
}

func TestWalkBlockStmt(t *testing.T) {
	n := &BasicLit{}
	b := &BlockStmt{List: []Node{n}}
	walk(t, b, []Node{b, n, nil, nil})
}

func TestWalkEnum(t *testing.T) {
	e := &Enum{Body: []Node{&Ident{}}}
	walk(t, e, []Node{e, e.Body[0], nil, nil})
}

func TestWalkEnumField(t *testing.T) {
	e := &EnumField{}
	walk(t, e, []Node{e, e.Name, nil, nil})
}

func TestWalkFile(t *testing.T) {
	f := &File{Nodes: []Node{&Ident{}}}
	walk(t, f, []Node{f, f.Nodes[0], nil, nil})
}

func TestWalkImport(t *testing.T) {
	im := &Import{Modifiers: []*Ident{&Ident{}}}
	walk(t, im, []Node{im, im.Modifiers[0], nil, im.Path, nil, nil})
}

func TestWalkMapType(t *testing.T) {
	m := &MapType{}
	walk(t, m, []Node{m, m.Key, nil, m.Value, nil, nil})
}

func TestWalkMessage(t *testing.T) {
	m := &Message{Body: []Node{&Ident{}}}
	walk(t, m, []Node{m, m.Name, nil, m.Body[0], nil, nil})
}
