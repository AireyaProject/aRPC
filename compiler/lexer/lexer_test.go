package lexer

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
		// A test file
		namespace aireya.example;

		enum Status {
			ACTIVE = 1;
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

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{NAMESPACE, "namespace"},
		{IDENT, "aireya.example"},
		{SEMI, ";"},

		{ENUM, "enum"},
		{IDENT, "Status"},
		{LBRACE, "{"},
		{IDENT, "ACTIVE"},
		{ASSIGN, "="},
		{NUMBER, "1"},
		{SEMI, ";"},
		{RBRACE, "}"},

		{STRUCT, "struct"},
		{IDENT, "User"},
		{LBRACE, "{"},
		{NUMBER, "1"},
		{COLON, ":"},
		{IDENT, "int64"},
		{IDENT, "id"},
		{SEMI, ";"},
		{NUMBER, "2"},
		{COLON, ":"},
		{IDENT, "string"},
		{IDENT, "name"},
		{SEMI, ";"},
		{NUMBER, "3"},
		{COLON, ":"},
		{IDENT, "list"},
		{LANGLE, "<"},
		{IDENT, "string"},
		{RANGLE, ">"},
		{IDENT, "tags"},
		{SEMI, ";"},
		{RBRACE, "}"},

		{SERVICE, "service"},
		{IDENT, "UserService"},
		{LBRACE, "{"},
		
		{RPC, "rpc"},
		{IDENT, "GetUser"},
		{LPAREN, "("},
		{IDENT, "User"},
		{RPAREN, ")"},
		{ARROW, "->"},
		{LPAREN, "("},
		{IDENT, "User"},
		{RPAREN, ")"},
		{SEMI, ";"},

		{RPC, "rpc"},
		{IDENT, "StreamUsers"},
		{LPAREN, "("},
		{IDENT, "User"},
		{RPAREN, ")"},
		{ARROW, "->"},
		{STREAM, "stream"},
		{LPAREN, "("},
		{IDENT, "User"},
		{RPAREN, ")"},
		{SEMI, ";"},

		{RBRACE, "}"},
		{EOF, ""},
	}

	l := New(strings.NewReader(input))

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType.String(), tok.Type.String())
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
