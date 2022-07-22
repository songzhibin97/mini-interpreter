package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/songzhibin97/mini-interpreter/token"
)

func TestLexer_NextToken(t *testing.T) {
	l := NewLexer(` + - * / % & | ^ < > = ! ( ) [ ] { } , . ; : << >> &^ += -= *= /= %= &= |= ^= <<= >>= &^= && || <- ++ -- == != <= >= := ... abc  123 "abc" "abc cba" macro`)
	tests := []*token.Token{
		{Type: token.ADD, Value: "+"},
		{Type: token.SUB, Value: "-"},
		{Type: token.MUL, Value: "*"},
		{Type: token.QUO, Value: "/"},
		{Type: token.REM, Value: "%"},
		{Type: token.AND, Value: "&"},
		{Type: token.OR, Value: "|"},
		{Type: token.XOR, Value: "^"},
		{Type: token.LSS, Value: "<"},
		{Type: token.GTR, Value: ">"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.NOT, Value: "!"},
		{Type: token.LPAREN, Value: "("},
		{Type: token.RPAREN, Value: ")"},
		{Type: token.LBRACK, Value: "["},
		{Type: token.RBRACK, Value: "]"},
		{Type: token.LBRACE, Value: "{"},
		{Type: token.RBRACE, Value: "}"},
		{Type: token.COMMA, Value: ","},
		{Type: token.PERIOD, Value: "."},
		{Type: token.SEMICOLON, Value: ";"},
		{Type: token.COLON, Value: ":"},
		{Type: token.SHL, Value: "<<"},
		{Type: token.SHR, Value: ">>"},
		{Type: token.AND_NOT, Value: "&^"},
		{Type: token.ADD_ASSIGN, Value: "+="},
		{Type: token.SUB_ASSIGN, Value: "-="},
		{Type: token.MUL_ASSIGN, Value: "*="},
		{Type: token.QUO_ASSIGN, Value: "/="},
		{Type: token.REM_ASSIGN, Value: "%="},
		{Type: token.AND_ASSIGN, Value: "&="},
		{Type: token.OR_ASSIGN, Value: "|="},
		{Type: token.XOR_ASSIGN, Value: "^="},
		{Type: token.SHL_ASSIGN, Value: "<<="},
		{Type: token.SHR_ASSIGN, Value: ">>="},
		{Type: token.AND_NOT_ASSIGN, Value: "&^="},
		{Type: token.LAND, Value: "&&"},
		{Type: token.LOR, Value: "||"},
		{Type: token.ARROW, Value: "<-"},
		{Type: token.INC, Value: "++"},
		{Type: token.DEC, Value: "--"},
		{Type: token.EQL, Value: "=="},
		{Type: token.NEQ, Value: "!="},
		{Type: token.LEQ, Value: "<="},
		{Type: token.GEQ, Value: ">="},
		{Type: token.DEFINE, Value: ":="},
		{Type: token.ELLIPSIS, Value: "..."},
		{Type: token.IDENT, Value: "abc"},
		{Type: token.INT, Value: "123"},
		{Type: token.STRING, Value: "abc"},
		{Type: token.STRING, Value: "abc cba"},
		{Type: token.MACRO, Value: "macro"},
		{Type: token.EOF, Value: ""},
	}
	for _, tt := range tests {
		tk := l.NextToken()
		assert.Equal(t, tt.Type, tk.Type)
		assert.Equal(t, tt.Value, tk.Value)
	}

	l = NewLexer(`
		var a = 10;
	    func add (a int, b int) int {
			return a + b 
		}
		
		type X interface {}
	`)
	tests = []*token.Token{
		{Type: token.VAR, Value: "var"},
		{Type: token.IDENT, Value: "a"},
		{Type: token.ASSIGN, Value: "="},
		{Type: token.INT, Value: "10"},
		{Type: token.SEMICOLON, Value: ";"},
		{Type: token.FUNC, Value: "func"},
		{Type: token.IDENT, Value: "add"},
		{Type: token.LPAREN, Value: "("},
		{Type: token.IDENT, Value: "a"},
		{Type: token.IDENT, Value: "int"},
		{Type: token.COMMA, Value: ","},
		{Type: token.IDENT, Value: "b"},
		{Type: token.IDENT, Value: "int"},
		{Type: token.RPAREN, Value: ")"},
		{Type: token.IDENT, Value: "int"},
		{Type: token.LBRACE, Value: "{"},
		{Type: token.RETURN, Value: "return"},
		{Type: token.IDENT, Value: "a"},
		{Type: token.ADD, Value: "+"},
		{Type: token.IDENT, Value: "b"},
		{Type: token.RBRACE, Value: "}"},
		{Type: token.TYPE, Value: "type"},
		{Type: token.IDENT, Value: "X"},
		{Type: token.INTERFACE, Value: "interface"},
		{Type: token.LBRACE, Value: "{"},
		{Type: token.RBRACE, Value: "}"},
		{Type: token.EOF, Value: ""},
	}
	for _, tt := range tests {
		tk := l.NextToken()
		assert.Equal(t, tt.Type, tk.Type)
		assert.Equal(t, tt.Value, tk.Value)
	}
}
