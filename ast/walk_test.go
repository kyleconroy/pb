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

func TestEmpty(t *testing.T) {
	for _, n := range []Node{
		&BasicLit{},
		&EmptyStmt{},
		&EnumField{},
		&Expr{},
		&File{},
		&Ident{},
		&Import{},
		&MapType{},
		&MessageField{},
		&OneOf{},
		&Option{},
		&Package{},
		&RPC{},
		&Service{},
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
	e := &EnumField{Name: &Ident{}}
	walk(t, e, []Node{e, e.Name, nil, nil})
}

func TestWalkFile(t *testing.T) {
	f := &File{Nodes: []Node{&Ident{}}}
	walk(t, f, []Node{f, f.Nodes[0], nil, nil})
}

func TestWalkImport(t *testing.T) {
	im := &Import{Path: &BasicLit{}, Modifiers: []*Ident{&Ident{}}}
	walk(t, im, []Node{im, im.Modifiers[0], nil, im.Path, nil, nil})
}

func TestWalkMapType(t *testing.T) {
	m := &MapType{Key: &Ident{}, Value: &Ident{}}
	walk(t, m, []Node{m, m.Key, nil, m.Value, nil, nil})
}

func TestWalkMessageField(t *testing.T) {
	m := &MessageField{Name: &Ident{}, Number: &BasicLit{}, Repeated: &Ident{}, Type: &Ident{}}
	walk(t, m, []Node{m, m.Name, nil, m.Number, nil, m.Type, nil, m.Repeated, nil, nil})
}

func TestWalkOneOf(t *testing.T) {
	o := &OneOf{Name: &Ident{}, Body: []Node{&Ident{}}}
	walk(t, o, []Node{o, o.Name, nil, o.Body[0], nil, nil})
}

func TestWalkOption(t *testing.T) {
	o := &Option{Names: []*Ident{&Ident{}}, Constant: &BasicLit{}}
	walk(t, o, []Node{o, o.Names[0], nil, o.Constant, nil, nil})
}

func TestWalkRPC(t *testing.T) {
	r := &RPC{Name: &Ident{}, InType: &Ident{}, OutType: &Ident{}}
	walk(t, r, []Node{r, r.Name, nil, r.InType, nil, r.OutType, nil, nil})
}

func TestWalkService(t *testing.T) {
	s := &Service{Name: &Ident{}, Body: &BlockStmt{}}
	walk(t, s, []Node{s, s.Name, nil, s.Body, nil, nil})
}
