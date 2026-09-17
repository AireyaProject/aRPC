package lexer

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS

	// Literals
	IDENT  // main, User, req
	NUMBER // 1, 123
	STRING // "hello"

	// Operators and Delimiters
	ASSIGN // =
	COLON  // :
	SEMI   // ;
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	LANGLE // <
	RANGLE // >
	ARROW  // ->
	COMMA  // ,

	// Keywords
	NAMESPACE // namespace
	STRUCT    // struct
	ENUM      // enum
	SERVICE   // service
	RPC       // rpc
	STREAM    // stream
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	STRING: "STRING",

	ASSIGN: "=",
	COLON:  ":",
	SEMI:   ";",
	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",
	LANGLE: "<",
	RANGLE: ">",
	ARROW:  "->",
	COMMA:  ",",

	NAMESPACE: "namespace",
	STRUCT:    "struct",
	ENUM:      "enum",
	SERVICE:   "service",
	RPC:       "rpc",
	STREAM:    "stream",
}

func (tok TokenType) String() string {
	if tok >= 0 && int(tok) < len(tokens) {
		return tokens[tok]
	}
	return "UNKNOWN"
}

var keywords = map[string]TokenType{
	"namespace": NAMESPACE,
	"struct":    STRUCT,
	"enum":      ENUM,
	"service":   SERVICE,
	"rpc":       RPC,
	"stream":    STREAM,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, '%s', %d:%d}", t.Type, t.Literal, t.Line, t.Col)
}
