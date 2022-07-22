package lexer

import (
	"unicode"

	"github.com/songzhibin97/mini-interpreter/token"
)

type Lexer struct {
	pos   int    // 解析器当前解析到的位置
	ln    int    // input 长度
	input []rune // 解析器需要解析的字符串
}

// next
// @Description: 获取下一个字符,将其pos移动到下一位
// @receiver l
// @return rune
func (l *Lexer) next() rune {
	if l.pos >= l.ln {
		// 0 => EOF
		return 0
	}

	ret := l.input[l.pos]
	l.pos++
	return ret
}

// peek
// @Description: 获取下一个字符,但不移动pos
// @param l:
// @param offset: 偏移量
// @return rune
func (l *Lexer) peek(offset int) rune {
	if l.pos+offset >= l.ln {
		return 0
	}
	return l.input[l.pos+offset]
}

func isLetter(v rune, index int) bool {
	return unicode.IsLetter(v) || v == '_' || (index != 0 && unicode.IsDigit(v))
}

func isDigit(v rune) bool {
	return unicode.IsDigit(v)
}

func (l *Lexer) letter() string {
	pos := l.pos
	for ; l.pos < l.ln; l.pos++ {
		v := l.input[l.pos]
		if !isLetter(v, l.pos-pos+1) {
			break
		}
	}
	ret := l.input[pos:l.pos]
	return string(ret)
}

func (l *Lexer) digit() string {
	pos := l.pos
	for ; l.pos < l.ln; l.pos++ {
		v := l.input[l.pos]
		if !isDigit(v) {
			break
		}
	}
	ret := l.input[pos:l.pos]
	return string(ret)
}

func (l *Lexer) string() string {
	pos := l.pos
	for ; l.pos < l.ln; l.pos++ {
		v := l.input[l.pos]
		if v == '"' {
			break
		}
	}
	ret := l.input[pos:l.pos]
	return string(ret)
}

func (l *Lexer) skipInterference() {
	for ; l.pos < l.ln; l.pos++ {
		switch l.input[l.pos] {
		case ' ':
		case '\n':
		case '\r':
		case '\t':
		default:
			return
		}
	}
}

// NextToken
// @Description: 解析获取下一个有效的 Token
// @receiver l
// @return *token.Token
func (l *Lexer) NextToken() *token.Token {
	var tk *token.Token
	l.skipInterference()
	v := l.next()
	switch v {
	case 0:
		tk = token.NewToken(token.EOF, "")
	case '"':
		tk = token.NewToken(token.STRING, l.string())
		l.next()
	case '+':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.ADD_ASSIGN, "+=")
			l.next()
		case '+':
			tk = token.NewToken(token.INC, "++")
			l.next()
		default:
			tk = token.NewToken(token.ADD, "+")
		}
	case '-':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.SUB_ASSIGN, "-=")
			l.next()
		case '-':
			tk = token.NewToken(token.DEC, "--")
			l.next()
		default:
			tk = token.NewToken(token.SUB, "-")
		}
	case '*':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.MUL_ASSIGN, "*=")
			l.next()
		default:
			tk = token.NewToken(token.MUL, "*")
		}
	case '/':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.QUO_ASSIGN, "/=")
			l.next()
		default:
			tk = token.NewToken(token.QUO, "/")
		}
	case '%':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.REM_ASSIGN, "%=")
			l.next()
		default:
			tk = token.NewToken(token.REM, "%")
		}
	case '&':
		switch l.peek(0) {
		case '^':
			switch l.peek(1) {
			case '=':
				tk = token.NewToken(token.AND_NOT_ASSIGN, "&^=")
				l.next()
			default:
				tk = token.NewToken(token.AND_NOT, "&^")
			}
			l.next()
		case '=':
			tk = token.NewToken(token.AND_ASSIGN, "&=")
			l.next()
		case '&':
			tk = token.NewToken(token.LAND, "&&")
			l.next()
		default:
			tk = token.NewToken(token.AND, "&")
		}
	case '|':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.OR_ASSIGN, "|=")
			l.next()
		case '|':
			tk = token.NewToken(token.LOR, "||")
			l.next()
		default:
			tk = token.NewToken(token.OR, "|")
		}
	case '^':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.XOR_ASSIGN, "^=")
			l.next()
		default:
			tk = token.NewToken(token.XOR, "^")
		}
	case '<':
		switch l.peek(0) {
		case '<':
			switch l.peek(1) {
			case '=':
				tk = token.NewToken(token.SHL_ASSIGN, "<<=")
				l.next()
			default:
				tk = token.NewToken(token.SHL, "<<")
			}
			l.next()
		case '-':
			tk = token.NewToken(token.ARROW, "<-")
			l.next()
		case '=':
			tk = token.NewToken(token.LEQ, "<=")
			l.next()
		default:
			tk = token.NewToken(token.LSS, "<")
		}
	case '>':
		switch l.peek(0) {
		case '>':
			switch l.peek(1) {
			case '=':
				tk = token.NewToken(token.SHR_ASSIGN, ">>=")
				l.next()
			default:
				tk = token.NewToken(token.SHR, ">>")
			}
			l.next()
		case '=':
			tk = token.NewToken(token.GEQ, ">=")
			l.next()
		default:
			tk = token.NewToken(token.GTR, ">")
		}
	case '=':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.EQL, "==")
			l.next()
		default:
			tk = token.NewToken(token.ASSIGN, "=")
		}
	case '!':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.NEQ, "!=")
			l.next()
		default:
			tk = token.NewToken(token.NOT, "!")
		}
	case '(':
		tk = token.NewToken(token.LPAREN, "(")
	case ')':
		tk = token.NewToken(token.RPAREN, ")")
	case '[':
		tk = token.NewToken(token.LBRACK, "[")
	case ']':
		tk = token.NewToken(token.RBRACK, "]")
	case '{':
		tk = token.NewToken(token.LBRACE, "{")
	case '}':
		tk = token.NewToken(token.RBRACE, "}")
	case ',':
		tk = token.NewToken(token.COMMA, ",")
	case '.':
		switch l.peek(0) {
		case '.':
			switch l.peek(1) {
			case '.':
				tk = token.NewToken(token.ELLIPSIS, "...")
				l.next()
				l.next()
			}
		default:
			tk = token.NewToken(token.PERIOD, ".")
		}
	case ';':
		tk = token.NewToken(token.SEMICOLON, ";")
	case ':':
		switch l.peek(0) {
		case '=':
			tk = token.NewToken(token.DEFINE, ":=")
			l.next()
		default:
			tk = token.NewToken(token.COLON, ":")
		}
	default:
		switch {
		case isLetter(v, 0):
			identifier := string(v) + l.letter()
			tk = token.NewToken(token.Lookup(identifier), identifier)
		case isDigit(v):
			tk = token.NewToken(token.INT, string(v)+l.digit())
		default:
			tk = token.NewToken(token.ILLEGAL, "")
		}
	}
	return tk
}

// NewLexer
// @Description: 创建新词法解析器
// @param input:
// @return *Lexer
func NewLexer(input string) *Lexer {
	v := &Lexer{
		input: []rune(input),
	}
	v.ln = len(v.input)
	return v
}
