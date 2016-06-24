package ast

import "testing"

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
}

func TestWalkBasicLit(t *testing.T) {
	n := &BasicLit{}
	walk(t, n, []Node{n, nil})
}

func TestWalkBlockStmt(t *testing.T) {
	n := &BasicLit{}
	b := &BlockStmt{List: []Node{n}}
	walk(t, b, []Node{b, n, nil, nil})
}

func TestWalkImport(t *testing.T) {
	in := &Ident{}
	im := &Import{Modifiers: []*Ident{in}}
	walk(t, im, []Node{im, in, nil, im.Path, nil})
}
