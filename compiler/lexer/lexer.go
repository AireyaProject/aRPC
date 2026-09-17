package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Lexer struct {
	r    *bufio.Reader
	line int
	col  int
}

func New(r io.Reader) *Lexer {
	return &Lexer{
		r:    bufio.NewReader(r),
		line: 1,
		col:  0,
	}
}

func (l *Lexer) readRune() (rune, error) {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return 0, err
	}
	if ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
	return ch, nil
}

func (l *Lexer) unreadRune(ch rune) {
	if ch == '\n' {
		l.line--
		// Note: col tracking isn't perfectly reverseable without a stack, 
		// but it's okay for simple error reporting in this context.
	} else {
		l.col--
	}
	_ = l.r.UnreadRune()
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespaceAndComments()

	ch, err := l.readRune()
	if err != nil {
		if err == io.EOF {
			return Token{Type: EOF, Literal: "", Line: l.line, Col: l.col}
		}
		return Token{Type: ILLEGAL, Literal: err.Error(), Line: l.line, Col: l.col}
	}

	startLine, startCol := l.line, l.col

	switch ch {
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(ch)}
	case ':':
		tok = Token{Type: COLON, Literal: string(ch)}
	case ';':
		tok = Token{Type: SEMI, Literal: string(ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(ch)}
	case '<':
		tok = Token{Type: LANGLE, Literal: string(ch)}
	case '>':
		tok = Token{Type: RANGLE, Literal: string(ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(ch)}
	case '-':
		next, err := l.readRune()
		if err == nil && next == '>' {
			tok = Token{Type: ARROW, Literal: "->"}
		} else {
			if err == nil {
				l.unreadRune(next)
			}
			tok = Token{Type: ILLEGAL, Literal: string(ch)}
		}
	default:
		if isLetter(ch) {
			l.unreadRune(ch)
			literal := l.readIdentifier()
			tok = Token{Type: LookupIdent(literal), Literal: literal}
		} else if isDigit(ch) {
			l.unreadRune(ch)
			tok = Token{Type: NUMBER, Literal: l.readNumber()}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(ch)}
		}
	}

	tok.Line = startLine
	tok.Col = startCol
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		ch, err := l.readRune()
		if err != nil {
			return
		}

		if unicode.IsSpace(ch) {
			continue
		}

		if ch == '/' {
			next, err := l.readRune()
			if err == nil && next == '/' {
				// Line comment
				for {
					c, err := l.readRune()
					if err != nil || c == '\n' {
						break
					}
				}
				continue
			} else {
				if err == nil {
					l.unreadRune(next)
				}
				l.unreadRune(ch)
				return
			}
		}

		l.unreadRune(ch)
		break
	}
}

func (l *Lexer) readIdentifier() string {
	var buf bytes.Buffer
	for {
		ch, err := l.readRune()
		if err != nil {
			break
		}
		if isLetter(ch) || isDigit(ch) || ch == '.' || ch == '_' {
			buf.WriteRune(ch)
		} else {
			l.unreadRune(ch)
			break
		}
	}
	return buf.String()
}

func (l *Lexer) readNumber() string {
	var buf bytes.Buffer
	for {
		ch, err := l.readRune()
		if err != nil {
			break
		}
		if isDigit(ch) {
			buf.WriteRune(ch)
		} else {
			l.unreadRune(ch)
			break
		}
	}
	return buf.String()
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
