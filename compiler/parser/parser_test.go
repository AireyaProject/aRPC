package parser

import (
	"strings"
	"testing"
	"github.com/aireya/aireyac/lexer"
)

func TestParser(t *testing.T) {
	input := `
		namespace aireya.example;

		enum Status {
			ACTIVE = 1;
			INACTIVE = 2;
		}

		struct User {
			1: int64 id;
			2: string name;
			3: list<string> tags;
		}

		service UserService {
			rpc GetUser(User) -> (User);
			rpc StreamUsers(User) -> stream(User);
		}
	`
	l := lexer.New(strings.NewReader(input))
	p := New(l)
	file := p.ParseFile()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser has errors: %v", p.Errors())
	}

	if file.Namespace != "aireya.example" {
		t.Errorf("expected namespace aireya.example, got %s", file.Namespace)
	}

	if len(file.Enums) != 1 {
		t.Fatalf("expected 1 enum, got %d", len(file.Enums))
	}
	if file.Enums[0].Name != "Status" || len(file.Enums[0].Values) != 2 {
		t.Errorf("enum Status parsing failed")
	}

	if len(file.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(file.Structs))
	}
	if file.Structs[0].Name != "User" || len(file.Structs[0].Fields) != 3 {
		t.Errorf("struct User parsing failed")
	}
	if file.Structs[0].Fields[2].Type != "list<string>" {
		t.Errorf("expected list<string>, got %s", file.Structs[0].Fields[2].Type)
	}

	if len(file.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(file.Services))
	}
	
	rpc1 := file.Services[0].RPCs[0]
	if rpc1.Name != "GetUser" || rpc1.ReqStream || rpc1.RespStream {
		t.Errorf("GetUser parsed incorrectly")
	}

	rpc2 := file.Services[0].RPCs[1]
	if rpc2.Name != "StreamUsers" || rpc2.ReqStream || !rpc2.RespStream {
		t.Errorf("StreamUsers parsed incorrectly")
	}
}
