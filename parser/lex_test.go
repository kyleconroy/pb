package parser

import (
	"strings"
	"testing"

	goast "go/ast"

	"github.com/kyleconroy/pb/ast"
)

const protoSimple = `syntax = "proto3";

import public "other.proto";

option java_package = "com.example.foo";

message message {
}

enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}

message outer {
  option (my_option).a = true;
  message inner {   // Level 2
    int64 ival = 1;
  } 
  enum OtherEnum {
    option allow_alias = true;
    UNKNOWN = 0;
    STARTED = 1;
    RUNNING = 2 [(custom_option) = "hello world"];
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}

service Limits {
  rpc Get(LimitsGetReq) returns (AccountLimits) {}
  rpc Set(LimitsSetReq) returns (Empty) {}
}
`

func TestLexer(t *testing.T) {
	f, err := ParseFile(strings.NewReader(protoSimple), 0)
	if err != nil {
		t.Error(err)
	}

	goast.Print(nil, f)

	if f.Syntax != ast.Proto3 {
		t.Error("The syntax should be proto3")
	}
}

func TestError(t *testing.T) {
	_, err := ParseFile(strings.NewReader("foo"), 0)
	if err == nil {
		t.Error(err)
	}
}
