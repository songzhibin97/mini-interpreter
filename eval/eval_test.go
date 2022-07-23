package eval

import (
	"testing"

	"github.com/songzhibin97/mini-interpreter/lexer"
	"github.com/songzhibin97/mini-interpreter/object"
	"github.com/songzhibin97/mini-interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func Test_evalIntegerExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalStringerExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`"hello" + " " + "world"`, "hello world"},
	}
	for _, tt := range tests {
		testStringerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalStringExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`"hello"`, "hello"},
		{`"test"`, "test"},
	}
	for _, tt := range tests {
		testStringerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalBuiltinFuncExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{`len("hello")`, 5},
		{`len(5)`, "argument to `len` not supported, got INT"},
	}
	for _, tt := range tests {
		switch v := tt.expect.(type) {
		case int:
			testIntegerObj(t, testEval(tt.input), int64(v))
		case string:
			testError(t, testEval(tt.input), v)
		}
	}
}

func Test_evalArray(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	v := testEval(input)
	array, ok := v.(*object.Array)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(array.Elements), 3)
	testIntegerObj(t, array.Elements[0], 1)
	testIntegerObj(t, array.Elements[1], 4)
	testIntegerObj(t, array.Elements[2], 6)
}

func Test_evalArrayIndexExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"var i = 0 ; [1][i]",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"var a = [1, 2, 3] a[2]",
			3,
		},
		{
			"var a = [1, 2, 3] a[0] + a[1] + a[2]",
			6,
		},
		{
			"var a = [1, 2, 3] var i = a[0] a[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		switch v := tt.expect.(type) {
		case int:
			testIntegerObj(t, testEval(tt.input), int64(v))
		default:
			testNilObj(t, testEval(tt.input))
		}
	}
}

func Test_evalMapExpr(t *testing.T) {
	input := `var two = "two"
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`
	expr := testEval(input)
	mp, ok := expr.(*object.Map)
	assert.Equal(t, ok, true)
	expect := map[object.MapKey]int64{
		(&object.Stringer{Value: "one"}).MapKey():   1,
		(&object.Stringer{Value: "two"}).MapKey():   2,
		(&object.Stringer{Value: "three"}).MapKey(): 3,
		(&object.Integer{Value: 4}).MapKey():        4,
		(&object.Boolean{Value: true}).MapKey():     5,
		(&object.Boolean{Value: false}).MapKey():    6,
	}
	assert.Equal(t, len(mp.Elements), len(expect))
	for k, v := range expect {
		pair, ok := mp.Elements[k]
		assert.Equal(t, ok, true)
		testIntegerObj(t, pair.Value, v)
	}
}

func Test_evalMapIndexExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		switch v := tt.expect.(type) {
		case int:
			testIntegerObj(t, testEval(tt.input), int64(v))
		default:
			testNilObj(t, testEval(tt.input))
		}
	}
}

func Test_evalBooleanExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		testBooleanObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalNotOperatorExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		testBooleanObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalIfExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		switch vv := tt.expect.(type) {
		case int:
			testIntegerObj(t, testEval(tt.input), int64(vv))
		default:
			testNilObj(t, testEval(tt.input))
		}
	}
}

func Test_evalReturnExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"return 10", 10},
		{"return 2 * 5", 10},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (10 > 1) { return 10 }", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10
			}
			return 1
		}
				`, 10},
		{`
		func a(x) {
			return x
		}(10)
				`, 10},
		{`
		func a (x) {
			return x
		}
		a(10)		
				`, 10},
	}
	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalVarExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"var a = 5 a", 5},
		{"var a = 5 * 5 a", 25},
		{"var a = 5 var b = a b", 5},
		{"var a = 5 var b = a var c = a + b + 5 c", 15},
	}
	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalFuncExpr(t *testing.T) {
	input := "func a (x) { x + 2 }"
	eval := testEval(input)
	fn, ok := eval.(*object.Function)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(fn.Parameters), 1)
	assert.Equal(t, fn.Parameters[0].String(), "x")
	assert.Equal(t, fn.Name.String(), "a")
	assert.Equal(t, fn.Body.String(), "(x + 2)")
}

func Test_evalCallExpr(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"func a(x) { x }  a(5)", 5},
		{"func a(x) { return x } a(5)", 5},
		{"func double(x) { x * 2 } double(5)", 10},
		{"func add(x, y) { x + y } add(5, 5)", 10},
		{"func add(x, y) { x + y } add(5 + 5, add(5, 5))", 20},
		{"func add(x) { x }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expect)
	}
}

func Test_evalEnv(t *testing.T) {
	input := `
var first = 10
var second = 10
var third = 10
func add (first) {
  var second = 20

  first + second + third
}

add(20) + first + second	
	`
	testIntegerObj(t, testEval(input), 70)
}

func Test_evalQuote(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			`quote(1)`,
			`1`},
		{
			`quote(1 + 2)`,
			`(1 + 2)`},
		{
			`quote(foobar)`,
			`foobar`},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
	}
	for _, tt := range tests {
		q, ok := testEval(tt.input).(*object.Quote)
		assert.Equal(t, ok, true)
		assert.NotNil(t, q.Node)
		assert.Equal(t, q.Node.String(), tt.expect)
	}
}

func Test_evalUnQuote(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			`quote(unquote(4))`, `4`,
		},
		{
			`quote(unquote(4 + 4))`,
			`8`,
		},
		{
			`quote(8 + unquote(4 + 4))`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(4 + 4) + 8)`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(true))`,
			`true`,
		},
		{
			`quote(unquote(true == false))`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`(4 + 4)`,
		},
		{
			`var quotedInfixExpression = quote(4 + 4) quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},
	}
	for _, tt := range tests {
		q, ok := testEval(tt.input).(*object.Quote)
		assert.Equal(t, ok, true)
		assert.NotNil(t, q.Node)
		assert.Equal(t, q.Node.String(), tt.expect)
	}
}

func TestDefinedMacro(t *testing.T) {
	input := `
   var n = 1
   func a (x, y) { x + y }
   macro mc (x, y) { x + y }
   `
	p := parser.NewParser(lexer.NewLexer(input))
	env := object.NewEnv(nil)
	program := p.ParseProgram()
	DefinedMacro(program, env)

	assert.Equal(t, len(program.Stmts), 2)
	_, ok := env.Get("n")
	assert.Equal(t, ok, false)
	_, ok = env.Get("a")
	assert.Equal(t, ok, false)
	obj, ok := env.Get("mc")
	assert.Equal(t, ok, true)
	macro, ok := obj.(*object.Macro)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(macro.Parameters), 2)
	assert.Equal(t, macro.Parameters[0].String(), "x")
	assert.Equal(t, macro.Parameters[1].String(), "y")
	assert.Equal(t, macro.Body.String(), "(x + y)")
}

func TestExpendMacro(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			`
			macro infixExpr() { quote(1 + 2) }
			infixExpr()
		`,
			`(1 + 2)`,
		},
		{
			`
			macro reverse(a, b) { quote(unquote(b) - unquote(a)) }
			reverse(2 + 2, 10 - 5)
`,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
			macro unless (condition, consequence, alternative) {
                     quote(if (!(unquote(condition))) {
                         unquote(consequence)
                     } else {
                         unquote(alternative)
}) }
unless(10 > 5, puts("not greater"), puts("greater"))
`,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}
	for _, tt := range tests {
		p := parser.NewParser(lexer.NewLexer(tt.input))
		env := object.NewEnv(nil)
		program := p.ParseProgram()
		DefinedMacro(program, env)
		expand := ExpandMacro(program, env)
		assert.Equal(t, expand.String(), parser.NewParser(lexer.NewLexer(tt.expect)).ParseProgram().String())
	}
}

func testEval(input string) object.Object {
	p := parser.NewParser(lexer.NewLexer(input))
	return Eval(p.ParseProgram(), object.NewEnv(nil))
}

func testIntegerObj(t *testing.T, obj object.Object, expect int64) {
	result, ok := obj.(*object.Integer)
	assert.Equal(t, ok, true)
	assert.Equal(t, result.Value, expect)
}

func testError(t *testing.T, obj object.Object, expect string) {
	result, ok := obj.(*object.Error)
	assert.Equal(t, ok, true)
	assert.Equal(t, result.Error, expect)
}

func testStringerObj(t *testing.T, obj object.Object, expect string) {
	result, ok := obj.(*object.Stringer)
	assert.Equal(t, ok, true)
	assert.Equal(t, result.Value, expect)
}

func testBooleanObj(t *testing.T, obj object.Object, expect bool) {
	result, ok := obj.(*object.Boolean)
	assert.Equal(t, ok, true)
	assert.Equal(t, result.Value, expect)
}

func testNilObj(t *testing.T, obj object.Object) {
	_, ok := obj.(*object.Nil)
	assert.Equal(t, ok, true)
}
