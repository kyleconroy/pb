package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyleconroy/pb/ast"
	"github.com/kyleconroy/pb/token"
)

func TestProtos(t *testing.T) {

	files, _ := ioutil.ReadDir("./_protos")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		t.Run(f.Name(), func(t *testing.T) {
			fset := token.NewFileSet()
			handle, err := os.Open(filepath.Join(".", "_protos", f.Name()))
			if err != nil {
				t.Error(err)
				return
			}
			f, err := ParseFile(fset, f.Name(), handle, 0)
			if err != nil {
				t.Error(err)
			}
			if f.Syntax != ast.Proto3 {
				t.Error("The syntax should be proto3")
			}
		})
	}
}

func TestError(t *testing.T) {
	fset := token.NewFileSet()
	_, err := ParseFile(fset, "", strings.NewReader("foo"), 0)
	if err == nil {
		t.Error(err)
	}
}
