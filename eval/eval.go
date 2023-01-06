package eval

import (
	"fmt"

	"github.com/songzhibin97/mini-interpreter/token"

	"github.com/songzhibin97/mini-interpreter/ast"
	"github.com/songzhibin97/mini-interpreter/object"
)

type Handler func(node ast.Node, env *object.Env) object.Object

func Eval(node ast.Node, env *object.Env, handler ...Handler) object.Object {
	handler = append(handler, defaultEval)
	return handler[0](node, env)
}

// default

func defaultEval(node ast.Node, env *object.Env) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.BlockStmt:
		return evalBlockStmt(n, env)

	case *ast.ExprStmt:
		return defaultEval(n.Expr, env)

	case *ast.ReturnStmt:
		ret := defaultEval(n.Value, env)
		if isError(ret) {
			return ret
		}
		return &object.Return{Value: ret}

	case *ast.VarStmt:
		ret := defaultEval(n.Value, env)
		if isError(ret) {
			return ret
		}
		env.Set(n.Name.Value, ret)

	case *ast.PrefixExpr:
		right := defaultEval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(n.Operator, right)

	case *ast.InfixExpr:
		left := defaultEval(n.Left, env)
		if isError(left) {
			return left
		}

		right := defaultEval(n.Right, env)
		if isError(left) {
			return left
		}
		return evalInfixExpr(n.Operator, left, right)

	case *ast.IfExpr:
		return evalIfExpr(n, env)

	case *ast.FuncExpr:
		obj := &object.Function{
			Name:       n.Name,
			Parameters: n.Params,
			Body:       n.Body,
			Env:        env,
		}
		env.Set(n.Name.Value, obj)
		return obj

	case *ast.CallExpr:
		if n.Func.TokenValue() == "quote" {
			return quote(n.Args[0], env)
		}

		fn := defaultEval(n.Func, env)
		if isError(fn) {
			return fn
		}

		args := evalExpr(n.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return callFunc(fn, args)

	case *ast.IndexExpr:
		left := defaultEval(n.Left, env)
		if isError(left) {
			return left
		}
		index := defaultEval(n.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpr(left, index)

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.Integer:
		return &object.Integer{Value: n.Value}

	case *ast.String:
		return &object.Stringer{Value: n.Value}

	case *ast.Boolean:
		return &object.Boolean{Value: n.Value}

	case *ast.Array:
		elements := evalExpr(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.Map:
		return evalMapExpr(n, env)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var r object.Object
	for _, stmt := range program.Stmts {
		r = Eval(stmt, env)

		switch r := r.(type) {
		case *object.Return:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return r
}

func evalBlockStmt(block *ast.BlockStmt, env *object.Env) object.Object {
	var r object.Object
	for _, stmt := range block.Stmts {
		r = Eval(stmt, env)
		if r == nil {
			continue
		}
		switch r.Type() {
		case object.RETURN, object.ERROR:
			return r
		}
	}
	return r
}

func evalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperatorExpr(right)
	case "-":
		return evalSubOperatorExpr(right)

	default:
		return &object.Error{Error: fmt.Sprintf("unknown prefix operator: " + operator + right.Type().String())}
	}
}

func evalNotOperatorExpr(right object.Object) object.Object {
	switch v := right.(type) {
	case *object.Boolean:
		return &object.Boolean{Value: !v.Value}
	case *object.Nil:
		return &object.Boolean{Value: true}
	default:
		return &object.Boolean{Value: false}
	}
}

func evalSubOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.INT {
		return &object.Error{Error: fmt.Sprintf("unknown sub operator: " + right.Type().String())}
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT && right.Type() == object.INT:
		return evalIntegerInfixExpr(operator, left, right)
	case left.Type() == object.String && right.Type() == object.String:
		return evalStringerInfixExpr(operator, left, right)
	case operator == "==" && left.Type() == right.Type():
		return &object.Boolean{Value: left.Inspect() == right.Inspect()}
	case operator == "!=" && left.Type() == right.Type():
		return &object.Boolean{Value: left.Inspect() != right.Inspect()}
	case left.Type() != right.Type():
		return &object.Error{Error: fmt.Sprintf("type mismatch: " + left.Type().String() + " " + operator + " " + right.Type().String())}
	default:
		return &object.Error{Error: fmt.Sprintf("unknown infix operator: " + left.Type().String() + " " + operator + " " + right.Type().String())}
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	builtin, ok := builtins[node.Value]
	if ok {
		return builtin
	}

	return val
}

func evalIfExpr(node *ast.IfExpr, env *object.Env) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return &object.Nil{}
	}
}

func evalMapExpr(node *ast.Map, env *object.Env) object.Object {
	elements := make(map[object.MapKey]object.HashValue)
	for k, v := range node.Elements {
		key := Eval(k, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.HashAble)
		if !ok {
			return &object.Error{Error: fmt.Sprintf("unable to hash key: " + key.Type().String())}
		}

		value := Eval(v, env)
		if isError(value) {
			return value
		}
		elements[hashKey.MapKey()] = object.HashValue{Key: key, Value: value}
	}
	return &object.Map{Elements: elements}
}

func isTruthy(obj object.Object) bool {
	switch v := obj.(type) {
	case *object.Nil:
		return false
	case *object.Boolean:
		return v.Value
	default:
		return true
	}
}

func evalExpr(args []ast.Expr, env *object.Env) []object.Object {
	var r []object.Object
	for _, arg := range args {
		eval := Eval(arg, env)
		if isError(eval) {
			return []object.Object{eval}
		}
		r = append(r, eval)
	}
	return r
}

func callFunc(fn object.Object, args []object.Object) object.Object {
	switch _fn := fn.(type) {
	case *object.Function:
		eval := Eval(_fn.Body, extendFuncEnv(_fn, args))
		return unwrapReturnValue(eval)
	case *object.Builtin:
		return _fn.Fn(args...)
	}
	return &object.Error{Error: fmt.Sprintf("not a function")}
}

func quote(node ast.Node, env *object.Env, modify ...ast.Modify) object.Object {
	modify = append(modify, ast.DefaultModify)
	return &object.Quote{Node: evalUnquoteCall(node, env, modify[0])}
}

func evalUnquoteCall(quote ast.Node, env *object.Env, modify ast.Modify) ast.Node {
	return modify(quote, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		callExpr := node.(*ast.CallExpr)
		if len(callExpr.Args) != 1 {
			return node
		}
		return translationObjToNode(Eval(callExpr.Args[0], env))
	})
}
func isUnquoteCall(quote ast.Node) bool {
	callExpr, ok := quote.(*ast.CallExpr)
	if !ok {
		return false
	}
	return callExpr.Func.TokenValue() == "unquote"
}

func translationObjToNode(obj object.Object) ast.Node {
	switch v := obj.(type) {
	case *object.Integer:
		return ast.Integer{
			Token: &token.Token{
				Type:  token.INT,
				Value: fmt.Sprintf("%d", v.Value),
			},
			Value: v.Value,
		}
	case *object.Boolean:
		t := &token.Token{
			Type:  token.FALSE,
			Value: "false",
		}
		if v.Value {
			t = &token.Token{
				Type:  token.TRUE,
				Value: "true",
			}
		}
		return ast.Boolean{
			Token: t,
			Value: v.Value,
		}

	case *object.Quote:
		return v.Node

	default:
		return nil
	}
}

func extendFuncEnv(fn *object.Function, args []object.Object) *object.Env {
	env := object.NewEnv(fn.Env)

	for index, parameter := range fn.Parameters {
		env.Set(parameter.Value, args[index])
	}
	return env
}

func evalIntegerInfixExpr(operator string, left, right object.Object) object.Object {
	l, r := left.(*object.Integer).Value, right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case "<":
		return &object.Boolean{Value: l < r}
	case ">":
		return &object.Boolean{Value: l > r}
	case "==":
		return &object.Boolean{Value: l == r}
	case "!=":
		return &object.Boolean{Value: l != r}
	default:
		return &object.Error{Error: fmt.Sprintf("unknown operator: " + operator + left.Type().String() + right.Type().String())}
	}
}

func evalStringerInfixExpr(operator string, left, right object.Object) object.Object {
	l, r := left.(*object.Stringer).Value, right.(*object.Stringer).Value
	switch operator {
	case "+":
		return &object.Stringer{Value: l + r}
	default:
		return &object.Error{Error: fmt.Sprintf("unknown operator: " + operator + left.Type().String() + right.Type().String())}
	}
}

func evalIndexExpr(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INT:
		return evalArrayIndexExpr(left, index)
	case left.Type() == object.MAP:
		return evalMapIndexExpr(left, index)
	default:
		return &object.Error{Error: fmt.Sprintf("index operator not supported: %s", left.Type().String())}
	}
}

func evalArrayIndexExpr(left object.Object, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(array.Elements) - 1)
	if idx < 0 || idx > max {
		return &object.Nil{}
	}
	return array.Elements[idx]
}

func evalMapIndexExpr(left object.Object, index object.Object) object.Object {
	mp := left.(*object.Map)
	key, ok := index.(object.HashAble)
	if !ok {
		return &object.Error{Error: fmt.Sprintf("unhashable type: %s", index.Type().String())}
	}
	val, ok := mp.Elements[key.MapKey()]
	if !ok {
		return &object.Nil{}
	}
	return val.Value
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.ERROR
}

func unwrapReturnValue(obj object.Object) object.Object {
	v, ok := obj.(*object.Return)
	if ok {
		return v.Value
	}
	return obj
}

func DefinedMacro(program *ast.Program, env *object.Env) {
	var definitions []int
	for index, stmt := range program.Stmts {
		if !isMacroDefinition(stmt) {
			continue
		}
		addMacro(stmt, env)
		definitions = append(definitions, index)
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		program.Stmts = append(program.Stmts[:definitions[i]], program.Stmts[definitions[i]+1:]...)
	}
}

func isMacroDefinition(node ast.Stmt) bool {
	stmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return false
	}
	_, ok = stmt.Expr.(*ast.Macro)
	return ok
}

func addMacro(stmt ast.Node, env *object.Env) {
	exprStmt, _ := stmt.(*ast.ExprStmt)
	macro, _ := exprStmt.Expr.(*ast.Macro)
	env.Set(macro.Name.Value, &object.Macro{
		Parameters: macro.Params,
		Body:       macro.Body,
		Env:        env,
	})
}

func ExpandMacro(program *ast.Program, env *object.Env, modify ...ast.Modify) ast.Node {
	modify = append(modify, ast.DefaultModify)
	return modify[0](program, func(node ast.Node) ast.Node {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExpr, env)
		if !ok {
			return node
		}
		args := quotaArgs(callExpr)
		evalEnv := extendMacroEnv(macro, args)
		eval := Eval(macro.Body, evalEnv)

		quote, ok := eval.(*object.Quote)
		if !ok {
			panic(fmt.Sprintf("invalid macro return value: %s", eval.Type().String()))
		}
		return quote.Node
	})
}

func quotaArgs(expr *ast.CallExpr) []*object.Quote {
	args := make([]*object.Quote, 0, len(expr.Args))
	for _, arg := range expr.Args {
		args = append(args, &object.Quote{Node: arg})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Env {
	extend := object.NewEnv(macro.Env)
	for i, parameter := range macro.Parameters {
		extend.Set(parameter.Value, args[i])
	}
	return extend
}

func isMacroCall(expr *ast.CallExpr, env *object.Env) (*object.Macro, bool) {
	ident, ok := expr.Func.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	return macro, true
}
