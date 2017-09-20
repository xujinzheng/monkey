package lexer

import (
	"github.com/xujinzheng/monkey/token"
)

type Lexer interface {
	NextToken() token.Token
}

type MonkeyLexer struct {
	input        string
	position     int  // current ch position
	readPosition int  // next ch position
	ch           byte // chrrent char
}

func NewMonkeyLexer(input string) *MonkeyLexer {
	l := &MonkeyLexer{
		input: input,
	}

	l.readChar()

	return l
}

func (p *MonkeyLexer) NextToken() token.Token {
	var tok token.Token

	p.skipWhitespace()

	switch p.ch {
	case '=':
		{
			if p.peekChar() == '=' { // ==
				ch := p.ch // current char
				p.readChar()
				tok = newStrToken(token.EQ, string(ch)+string(p.ch))
			} else {
				tok = newToken(token.ASSIGN, p.ch)
			}
		}
	case '+', '-', '*', '/', '<', '>', ';', ',', '{', '}', '(', ')':
		{
			tokenType := charToTokenType[p.ch]
			tok = newToken(tokenType, p.ch)
		}
	case '!':
		{
			if p.peekChar() == '=' { // ==
				ch := p.ch // current char
				p.readChar()
				tok = newStrToken(token.NOT_EQ, string(ch)+string(p.ch))
			} else {
				tok = newToken(token.BANG, p.ch)
			}
		}
	case 0:
		{
			tok.Literal = ""
			tok.Type = token.EOF
		}
	default:
		{
			if isLetter(p.ch) {
				identStr := p.readIdentifer()
				typ := token.LookupIdent(identStr) // keywords
				tok = newStrToken(typ, identStr)
				return tok
			} else if isDigit(p.ch) {
				num := p.readNumber()
				tok = newStrToken(token.INT, num)
				return tok
			} else {
				tok = newToken(token.ILLEGAL, p.ch) // we do not recognize this char ....
			}
		}
	}

	p.readChar()

	return tok
}

func (p *MonkeyLexer) readChar() {
	if p.readPosition >= len(p.input) {
		p.ch = 0
		return
	}

	p.ch = p.input[p.readPosition]

	p.position = p.readPosition
	p.readPosition++
}

func (p *MonkeyLexer) peekChar() byte {
	if p.readPosition >= len(p.input) {
		p.ch = 0
		return 0
	}

	return p.input[p.readPosition]
}

func (p *MonkeyLexer) skipWhitespace() {
	for p.ch == ' ' || p.ch == '\t' || p.ch == '\n' || p.ch == '\r' {
		p.readChar()
	}
}

func (p *MonkeyLexer) readNumber() string {
	position := p.position //backup start position

	for isDigit(p.ch) {
		p.readChar()
	}

	return p.input[position:p.position]
}

func (p *MonkeyLexer) readIdentifer() string {

	position := p.position //backup start position

	for isLetter(p.ch) {
		p.readChar()
	}

	return p.input[position:p.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newStrToken(tokenType token.TokenType, ch string) token.Token {
	return token.Token{Type: tokenType, Literal: ch}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

var (
	charToTokenType = map[byte]token.TokenType{
		'+': token.PLUS,
		'-': token.MINUS,
		'*': token.ASTERISK,
		'/': token.SLASH,
		'<': token.LT,
		'>': token.GT,
		';': token.SEMICOLON,
		',': token.COMMA,
		'{': token.LBRACE,
		'}': token.RBRACE,
		'(': token.LPAREN,
		')': token.RPAREN,
	}
)
