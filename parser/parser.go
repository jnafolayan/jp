package parser

import (
	"fmt"

	"github.com/jnafolayan/json-parser/elements"
	"github.com/jnafolayan/json-parser/lexer"
	"github.com/jnafolayan/json-parser/tokens"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	return p
}

// Parse returns `true` is the JSON is valid and `false` otherwise.
// Parses the grammar found at https://www.json.org/json-en.html
func (p *Parser) Parse() error {
	_, err := p.parseElement()
	return err
}

func (p *Parser) nextToken() tokens.Token {
	p.currentToken = p.lexer.NextToken()
	return p.currentToken
}

func (p *Parser) parseElement() (elements.Element, error) {
	switch p.currentToken.Type {
	case tokens.LBRACE:
		return p.parseObject()
	case tokens.STRING:
		return p.parseString()
	case tokens.TRUE:
		fallthrough
	case tokens.FALSE:
		fallthrough
	case tokens.NULL:
		return p.parseKeyword()
	default:
		return nil, fmt.Errorf("unsupported token %q", p.currentToken.Literal)
	}
}

func (p *Parser) parseObject() (*elements.Object, error) {
	obj := &elements.Object{}

	p.nextToken()
	for p.currentToken.Type == tokens.STRING {
		key := p.currentToken
		p.nextToken()
		if p.currentToken.Type != tokens.COLON {
			return nil, fmt.Errorf("expected %q, got %q", tokens.COLON, p.currentToken.Literal)
		}

		p.nextToken()
		element, err := p.parseElement()
		if err != nil {
			return nil, err
		}

		obj.Pairs = append(obj.Pairs, &elements.ObjectPair{
			Key:   key,
			Value: element,
		})

		p.nextToken()
		if p.currentToken.Type == tokens.COMMA {
			p.nextToken()
		}
	}

	if p.currentToken.Type == tokens.RBRACE {
		return obj, nil
	}

	return nil, fmt.Errorf("expected %q, found %q", tokens.LBRACE, p.currentToken.Literal)
}

func (p *Parser) parseString() (*elements.String, error) {
	return &elements.String{Value: p.currentToken}, nil
}

func (p *Parser) parseKeyword() (*elements.Keyword, error) {
	return &elements.Keyword{Value: p.currentToken}, nil

}
