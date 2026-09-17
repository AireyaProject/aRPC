package parser

import (
	"fmt"
	"github.com/aireya/aireyac/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) ParseFile() *File {
	file := &File{}

	for p.curToken.Type != lexer.EOF {
		switch p.curToken.Type {
		case lexer.NAMESPACE:
			if p.expectPeek(lexer.IDENT) {
				file.Namespace = p.curToken.Literal
				if p.expectPeek(lexer.SEMI) {
					p.nextToken()
				}
			}
		case lexer.ENUM:
			file.Enums = append(file.Enums, p.parseEnum())
		case lexer.STRUCT:
			file.Structs = append(file.Structs, p.parseStruct())
		case lexer.SERVICE:
			file.Services = append(file.Services, p.parseService())
		default:
			// skip unknown for now or record error
			p.nextToken()
		}
	}
	return file
}

func (p *Parser) parseEnum() *Enum {
	e := &Enum{}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	e.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken() // advance past {

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENT {
			name := p.curToken.Literal
			if p.expectPeek(lexer.ASSIGN) && p.expectPeek(lexer.NUMBER) {
				val := p.curToken.Literal
				e.Values = append(e.Values, &EnumValue{Name: name, Value: val})
				p.expectPeek(lexer.SEMI) // optional check
			}
		}
		p.nextToken()
	}
	p.nextToken() // advance past }
	return e
}

func (p *Parser) parseStruct() *Struct {
	s := &Struct{}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	s.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken() // advance past {

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.NUMBER {
			f := &Field{Tag: p.curToken.Literal}
			if p.expectPeek(lexer.COLON) {
				p.nextToken() // advance to type
				f.Type = p.parseType()
				if p.expectPeek(lexer.IDENT) {
					f.Name = p.curToken.Literal
					s.Fields = append(s.Fields, f)
					p.expectPeek(lexer.SEMI)
				}
			}
		}
		p.nextToken()
	}
	p.nextToken() // advance past }
	return s
}

func (p *Parser) parseType() string {
	typ := p.curToken.Literal
	if p.peekToken.Type == lexer.LANGLE {
		p.nextToken() // advance to <
		typ += "<"
		p.nextToken() // advance to inner type
		typ += p.parseType()
		if p.expectPeek(lexer.RANGLE) {
			typ += ">"
		}
	}
	return typ
}

func (p *Parser) parseService() *Service {
	s := &Service{}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	s.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken() // advance past {

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.RPC {
			s.RPCs = append(s.RPCs, p.parseRPC())
		}
		p.nextToken()
	}
	p.nextToken() // advance past }
	return s
}

func (p *Parser) parseRPC() *RPC {
	r := &RPC{}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	r.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken() // advance into parens

	if p.curToken.Type == lexer.STREAM {
		r.ReqStream = true
		p.nextToken()
	}
	if p.curToken.Type == lexer.IDENT {
		r.ReqType = p.curToken.Literal
	}
	p.expectPeek(lexer.RPAREN)

	if !p.expectPeek(lexer.ARROW) {
		return nil
	}

	// Response part
	if p.peekToken.Type == lexer.STREAM {
		p.nextToken() // consume stream
		r.RespStream = true
		if p.expectPeek(lexer.LPAREN) {
			p.nextToken() // into parens
			if p.curToken.Type == lexer.IDENT {
				r.RespType = p.curToken.Literal
			}
			p.expectPeek(lexer.RPAREN)
		}
	} else if p.expectPeek(lexer.LPAREN) {
		p.nextToken() // into parens
		if p.curToken.Type == lexer.IDENT {
			r.RespType = p.curToken.Literal
		}
		p.expectPeek(lexer.RPAREN)
	}

	p.expectPeek(lexer.SEMI)
	return r
}
