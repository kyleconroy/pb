package parser

import "testing"

const protoSimple = `syntax = "proto3";

import public "other.proto";
option java_package = "com.example.foo";
package foo.bar;

message SearchRequest {
  string query = 1; // Hello there main
  int32 page_number = 2;
  int32 result_per_page = 3;
}
`

func TestLexer(t *testing.T) {
	_, items := lex("protoSimple", protoSimple)
	for i := range items {
		t.Logf("%s", i)
	}
}

func TestError(t *testing.T) {
	_, items := lex("protoError", "foo")
	for i := range items {
		t.Logf("%s", i)
	}
}
