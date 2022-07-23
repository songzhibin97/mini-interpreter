package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/songzhibin97/mini-interpreter/ast"
	"github.com/songzhibin97/mini-interpreter/lexer"
	"github.com/stretchr/testify/assert"
)

// ============================================

func testVarStmt(t *testing.T, s ast.Stmt, name string) {
	assert.Equal(t, s.TokenValue(), "var")
	varStmt, ok := s.(*ast.VarStmt)
	assert.Equal(t, ok, true)
	assert.Equal(t, varStmt.Name.Value, name)
	assert.Equal(t, varStmt.Name.TokenValue(), name)
}

func testInteger(t *testing.T, expr ast.Expr, value int64) {
	integer, ok := expr.(*ast.Integer)
	assert.Equal(t, ok, true)
	assert.Equal(t, integer.Value, value)
	assert.Equal(t, integer.TokenValue(), strconv.Itoa(int(value)))
}

func testIdentifier(t *testing.T, expr ast.Expr, value string) {
	ident, ok := expr.(*ast.Identifier)
	assert.Equal(t, ok, true)
	assert.Equal(t, ident.Value, value)
	assert.Equal(t, ident.TokenValue(), value)
}

func testBoolean(t *testing.T, expr ast.Expr, value bool) {
	b, ok := expr.(*ast.Boolean)
	assert.Equal(t, ok, true)
	assert.Equal(t, b.Value, value)
	assert.Equal(t, b.TokenValue(), fmt.Sprintf("%t", value))
}

func testInfixExpr(t *testing.T, expr ast.Expr, left interface{}, op string, right interface{}) {
	opExpr, ok := expr.(*ast.InfixExpr)
	assert.Equal(t, ok, true)
	testExpr(t, opExpr.Left, left)
	assert.Equal(t, opExpr.Operator, op)
	testExpr(t, opExpr.Right, right)
}

func testExpr(t *testing.T, expr ast.Expr, expect interface{}) {
	switch v := expect.(type) {
	case int:
		testInteger(t, expr, int64(v))
	case int64:
		testInteger(t, expr, v)
	case string:
		testIdentifier(t, expr, v)
	case bool:
		testBoolean(t, expr, v)
	default:
		t.Errorf("type of exp not handled. got=%T", expect)
	}
}

// ============================================
func TestParser_parseVarStmt(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		value      interface{}
	}{
		{"var a = 1", "a", 1},
		{"var b = test", "b", "test"},
	}

	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt := v.Stmts[0]
		testVarStmt(t, stmt, tt.identifier)
	}
}

func TestParser_parseReturnStmt(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{"return 10", 10},
		{"return true", true},
		{"return a", "a"},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt, ok := v.Stmts[0].(*ast.ReturnStmt)
		assert.Equal(t, ok, true)
		assert.Equal(t, stmt.TokenValue(), "return")
		testExpr(t, stmt.Value, tt.expect)
	}
}

func TestParser_parseIdentifier(t *testing.T) {
	input := `test`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	identifier, ok := stmt.Expr.(*ast.Identifier)
	assert.Equal(t, ok, true)
	assert.Equal(t, identifier.Value, "test")
	assert.Equal(t, identifier.TokenValue(), "test")

}

func TestParser_parseInteger(t *testing.T) {
	input := `10`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	testInteger(t, stmt.Expr, int64(10))
}

func TestParser_parseString(t *testing.T) {
	input := `"hello"`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	integer, ok := stmt.Expr.(*ast.String)
	assert.Equal(t, ok, true)
	assert.Equal(t, integer.Value, "hello")
}

func TestParser_parseArray(t *testing.T) {
	input := `[]`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	integer, ok := stmt.Expr.(*ast.Array)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(integer.Elements), 0)
}

func TestParser_parseArrayElements(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	integer, ok := stmt.Expr.(*ast.Array)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(integer.Elements), 3)
	testExpr(t, integer.Elements[0], 1)
	testInfixExpr(t, integer.Elements[1], 2, "*", 2)
	testInfixExpr(t, integer.Elements[2], 3, "+", 3)
}

func TestParser_parseIndexExpr(t *testing.T) {
	input := `a[1+1]`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	integer, ok := stmt.Expr.(*ast.IndexExpr)
	assert.Equal(t, ok, true)
	testIdentifier(t, integer.Left, "a")
	testInfixExpr(t, integer.Index, 1, "+", 1)
}

func TestParser_parseMapEmpty(t *testing.T) {
	input := `{}`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	mp, ok := stmt.Expr.(*ast.Map)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(mp.Elements), 0)
}

func TestParser_parseMap(t *testing.T) {
	input := `{"a": 1, "b": 2, "c": 3}`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	mp, ok := stmt.Expr.(*ast.Map)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(mp.Elements), 3)
	expect := map[string]int64{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	for k, v := range mp.Elements {
		kk, ok := k.(*ast.String)
		assert.Equal(t, ok, true)
		testInteger(t, v, expect[kk.Value])
	}
}

func TestParser_parsePairExpr(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	mp, ok := stmt.Expr.(*ast.Map)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(mp.Elements), 3)
	tests := map[string]func(expr ast.Expr){
		"one": func(e ast.Expr) {
			testInfixExpr(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expr) {
			testInfixExpr(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expr) {
			testInfixExpr(t, e, 15, "/", 5)
		},
	}
	for k, v := range mp.Elements {
		kk, ok := k.(*ast.String)
		assert.Equal(t, ok, true)
		fn, ok := tests[kk.Value]
		assert.Equal(t, ok, true)
		fn(v)
	}
}

func TestParser_parsePrefixExpr(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!10", "!", 10},
		{"-10", "-", 10},
		{"!a", "!", "a"},
		{"-a", "-", "a"},
		{"!true", "!", true},
		{"!false", "!", false},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt, ok := v.Stmts[0].(*ast.ExprStmt)
		assert.Equal(t, ok, true)
		expr, ok := stmt.Expr.(*ast.PrefixExpr)
		assert.Equal(t, ok, true)
		assert.Equal(t, expr.Operator, tt.operator)
		testExpr(t, expr.Right, tt.value)
	}

}

func TestParser_MacroExpr(t *testing.T) {
	input := `macro t (x, y) { x + y }`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	mp, ok := stmt.Expr.(*ast.Macro)
	assert.Equal(t, ok, true)
	testIdentifier(t, mp.Name, "t")
	assert.Equal(t, len(mp.Params), 2)
	testIdentifier(t, mp.Params[0], "x")
	testIdentifier(t, mp.Params[1], "y")
	assert.Equal(t, len(mp.Body.Stmts), 1)
	testInfixExpr(t, mp.Body.Stmts[0].(*ast.ExprStmt).Expr, "x", "+", "y")
}

func TestParser_parseInfixExpr(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"1 + 1", 1, "+", 1},
		{"1 - 1", 1, "-", 1},
		{"1 * 1", 1, "*", 1},
		{"1 / 1", 1, "/", 1},
		{"1 > 1", 1, ">", 1},
		{"1 < 1", 1, "<", 1},
		{"1 == 1", 1, "==", 1},
		{"1 != 1", 1, "!=", 1},
		{"a + b", "a", "+", "b"},
		{"a - b", "a", "-", "b"},
		{"a * b", "a", "*", "b"},
		{"a / b", "a", "/", "b"},
		{"a > b", "a", ">", "b"},
		{"a < b", "a", "<", "b"},
		{"a == b", "a", "==", "b"},
		{"a != b", "a", "!=", "b"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range tests {
		p := NewParser(lexer.NewLexer(test.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt, ok := v.Stmts[0].(*ast.ExprStmt)
		assert.Equal(t, ok, true)
		expr, ok := stmt.Expr.(*ast.InfixExpr)
		assert.Equal(t, ok, true)
		testExpr(t, expr.Left, test.leftValue)
		assert.Equal(t, expr.Operator, test.operator)
		testExpr(t, expr.Right, test.rightValue)
	}
}

func TestParser_operator(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, v.String(), tt.expect)
	}
}

func TestParser_parseBooleanExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt, ok := v.Stmts[0].(*ast.ExprStmt)
		assert.Equal(t, ok, true)
		boolean, ok := stmt.Expr.(*ast.Boolean)
		assert.Equal(t, ok, true)
		assert.Equal(t, boolean.Value, tt.expect)
	}
}

func TestParser_parseIfExpr(t *testing.T) {
	input := `if (x < y) { x }`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	expr, ok := stmt.Expr.(*ast.IfExpr)
	assert.Equal(t, ok, true)
	testInfixExpr(t, expr.Condition, "x", "<", "y")
	assert.Equal(t, len(expr.Consequence.Stmts), 1)
	consequence, ok := expr.Consequence.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	testIdentifier(t, consequence.Expr, "x")
}

func TestParser_parseFuncExpr(t *testing.T) {
	input := `func a (x, y) { x + y }`
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	expr, ok := stmt.Expr.(*ast.FuncExpr)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(expr.Params), 2)
	testIdentifier(t, expr.Name, "a")
	testIdentifier(t, expr.Params[0], "x")
	testIdentifier(t, expr.Params[1], "y")
	assert.Equal(t, len(expr.Body.Stmts), 1)
	bodyStmt, ok := expr.Body.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	testInfixExpr(t, bodyStmt.Expr, "x", "+", "y")
}

func TestParser_parseFuncParams(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{"func a () {}", []string{}},
		{"func a (x) {}", []string{"x"}},
		{"func a (x, y, z) {}", []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		stmt, ok := v.Stmts[0].(*ast.ExprStmt)
		assert.Equal(t, ok, true)
		fn, ok := stmt.Expr.(*ast.FuncExpr)
		assert.Equal(t, ok, true)
		assert.Equal(t, len(fn.Params), len(tt.expect))
		for i, s := range tt.expect {
			testIdentifier(t, fn.Params[i], s)
		}
	}
}

func TestParser_parseCallExpr(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	p := NewParser(lexer.NewLexer(input))
	v := p.ParseProgram()
	for _, s := range p.Errors() {
		t.Errorf("parser error: %s", s)
	}
	assert.Equal(t, len(v.Stmts), 1)
	stmt, ok := v.Stmts[0].(*ast.ExprStmt)
	assert.Equal(t, ok, true)
	expr, ok := stmt.Expr.(*ast.CallExpr)
	assert.Equal(t, ok, true)
	testIdentifier(t, expr.Func, "add")
	assert.Equal(t, len(expr.Args), 3)
	testExpr(t, expr.Args[0], 1)
	testInfixExpr(t, expr.Args[1], 2, "*", 3)
	testInfixExpr(t, expr.Args[2], 4, "+", 5)
}

func TestParser_parseCallArgsExpr(t *testing.T) {
	tests := []struct {
		input string
		ident string
		args  []string
	}{
		{
			"add()",
			"add",
			[]string{},
		},
		{
			"add(1)",
			"add",
			[]string{"1"},
		},
		{
			"add(1, 2 * 3, 4 + 5)",
			"add",
			[]string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}
	for _, tt := range tests {
		p := NewParser(lexer.NewLexer(tt.input))
		v := p.ParseProgram()
		for _, s := range p.Errors() {
			t.Errorf("parser error: %s", s)
		}
		assert.Equal(t, len(v.Stmts), 1)
		stmt, ok := v.Stmts[0].(*ast.ExprStmt)
		assert.Equal(t, ok, true)
		expr, ok := stmt.Expr.(*ast.CallExpr)
		testIdentifier(t, expr.Func, tt.ident)
		assert.Equal(t, len(expr.Args), len(tt.args))
		for i, arg := range tt.args {
			assert.Equal(t, expr.Args[i].String(), arg)
		}
	}
}
