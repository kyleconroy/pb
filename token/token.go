package token

type FileSet struct {
}

type Pos int

type Token int

func NewFileSet() *FileSet {
	return &FileSet{}
}
