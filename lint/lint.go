package lint

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/kyleconroy/pb/ast"
	"github.com/kyleconroy/pb/parser"
	"github.com/kyleconroy/pb/token"
)

const styleGuideBase = "https://golang.org/wiki/CodeReviewComments"

// A Linter lints Go source code.
type Linter struct {
}

// Problem represents a problem in some source code.
type Problem struct {
	Position   token.Pos // position in source file
	Text       string    // the prose that describes the problem
	Link       string    // (optional) the link to the style guide for the problem
	Confidence float64   // a value in (0,1] estimating the confidence in this problem's correctness
	LineText   string    // the source line
	Category   string    // a short name for the general category of the problem
}

func Lint(filename string, src []byte) ([]Problem, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filename, bytes.NewBuffer(src), 0)
	if err != nil {
		return nil, err
	}
	h := file{
		f:        f,
		src:      src,
		filename: filename,
	}
	h.lint()
	return h.problems, nil
}

// file represents a protocol buffer file being linted.
type file struct {
	f        *ast.File
	src      []byte
	filename string
	problems []Problem
}

func (f *file) lint() {
	f.lintEnums()
	f.lintEnumFields()
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.f)
}

// walker adapts a function to satisfy the ast.Visitor interface.
// The function return whether the walk should proceed into the node's children.
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

// The variadic arguments may start with link and category types,
// and must end with a format string and any arguments.
// It returns the new Problem.
func (f *file) errorf(n ast.Node, confidence float64, msg string, args ...interface{}) {
	f.problems = append(f.problems, Problem{
		Position: n.Pos(),
		Text:     fmt.Sprintf(msg, args...),
	})
}

var camelCaseRE = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
var upperCaseRE = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
var snakeCaseRE = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// lintEnums complains if the name of an enum is not CamelCase.
func (f *file) lintEnums() {
	f.walk(func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.Enum:
			if v.Name != nil && !camelCaseRE.MatchString(v.Name.Name) {
				f.errorf(v, 0.9, "enum names should be CamelCase; %s", v.Name.Name)
			}
			return false
		}
		return true
	})
}

// lintEnumFields complains if the name of an enum field is not ALL_CAPS.
func (f *file) lintEnumFields() {
	f.walk(func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.EnumField:
			if v.Name != nil && !upperCaseRE.MatchString(v.Name.Name) {
				f.errorf(v, 0.9, "enum field names should be ALL_CAPS; %s", v.Name.Name)
			}
			return false
		}
		return true
	})
}
