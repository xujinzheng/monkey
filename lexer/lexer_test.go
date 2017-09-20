package lexer

import (
	"github.com/golang/mock/gomock"
	"os"
	"testing"

	"github.com/xujinzheng/monkey/lexer/mock_lexer"
	"github.com/xujinzheng/monkey/token"
)

func getTestData() (input string, tokens []token.Token) {
	input = `
let five = 5;
let ten = 10;
`

	tokens = []token.Token{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	return
}

func TestNextToken(t *testing.T) {

	input, excepts := getTestData()

	l, fn := newLexer(input, excepts, t)

	defer fn()

	for i, tt := range excepts {
		tok := l.NextToken()

		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.Type, tok.Type)
		}

		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.Literal, tok.Literal)
		}
	}
}

func newLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {
	env := os.Getenv("GO_MOCK_TEST")
	if env == "1" {
		t.Log("MOCK TEST ENABLED!!!")
		return newMockLexer(input, excepts, t)
	}

	return newMonkeyLexer(input, excepts, t)
}

func newMonkeyLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {
	return NewMonkeyLexer(input), func() {}
}

func newMockLexer(input string, excepts []token.Token, t *testing.T) (l Lexer, deferFN func()) {

	ctrl := gomock.NewController(t)

	mockLexter := mock_lexer.NewMockLexer(ctrl)

	for i := 0; i < len(excepts); i++ {
		mockLexter.EXPECT().NextToken().Return(excepts[i])
	}

	l = mockLexter
	deferFN = func() { ctrl.Finish() }

	return
}
