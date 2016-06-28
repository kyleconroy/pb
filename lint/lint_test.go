package lint

import "testing"

const sloppyEnum = `
syntax = "proto3";

package pb.lint;

enum UPPER_CASE {
  camelCase = 0;
  snake_case = 1;
  TitleCase = 2;
  UPPPER_CASE = 3;
}
`

func TestEnumLint(t *testing.T) {
	problems, err := Lint("", []byte(sloppyEnum))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(problems)
}
