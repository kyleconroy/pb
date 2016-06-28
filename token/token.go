package token

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	INT           // 12345
	STRING        // "abc"
	BOOL          // true | false
)
