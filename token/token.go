package token

type FileSet struct {
}

type Pos int

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	INT           // 12345
	STRING        // "abc"
	BOOL          // true | false
)

func NewFileSet() *FileSet {
	return &FileSet{}
}
