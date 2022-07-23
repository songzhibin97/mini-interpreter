package ast

import (
	"strings"

	"github.com/songzhibin97/mini-interpreter/token"
)

type Node interface {
	TokenValue() string
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenValue() string {
	if len(p.Stmts) != 0 {
		return p.Stmts[0].TokenValue()
	}
	return ""
}

func (p *Program) String() string {
	var b strings.Builder
	for _, stmt := range p.Stmts {
		b.WriteString(stmt.String())
	}
	return b.String()
}

// ============================================================================

// var <标识符> = <表达式>

type VarStmt struct {
	Token *token.Token
	Name  *Identifier
	Value Expr
}

func (v VarStmt) TokenValue() string { return v.Token.Value }
func (v VarStmt) stmtNode()          {}
func (v VarStmt) String() string {
	var b strings.Builder

	b.WriteString(v.TokenValue() + " ")
	b.WriteString(v.Name.String() + " = ")
	if v.Value != nil {
		b.WriteString(v.Value.String())
	}
	return b.String()
}

// ============================================================================

// return <表达式>

type ReturnStmt struct {
	Token *token.Token
	Value Expr
}

func (r ReturnStmt) TokenValue() string { return r.Token.Value }
func (r ReturnStmt) stmtNode()          {}
func (r ReturnStmt) String() string {
	var b strings.Builder
	b.WriteString(r.TokenValue() + " ")
	if r.Value != nil {
		b.WriteString(r.Value.String())
	}
	return b.String()
}

// ============================================================================

type ExprStmt struct {
	Token *token.Token
	Expr  Expr
}

func (e ExprStmt) TokenValue() string { return e.Token.Value }
func (e ExprStmt) stmtNode()          {}
func (e ExprStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

// ============================================================================

type BlockStmt struct {
	Token *token.Token
	Stmts []Stmt
}

func (b BlockStmt) TokenValue() string { return b.Token.Value }
func (b BlockStmt) stmtNode()          {}
func (b BlockStmt) String() string {
	var bb strings.Builder
	for _, stmt := range b.Stmts {
		bb.WriteString(stmt.String())
	}
	return bb.String()
}

// ============================================================================
// ============================================================================

type Identifier struct {
	Token *token.Token
	Value string
}

func (i Identifier) TokenValue() string { return i.Token.Value }
func (i Identifier) exprNode()          {}
func (i Identifier) String() string     { return i.Value }

// ============================================================================

type Boolean struct {
	Token *token.Token
	Value bool
}

func (b Boolean) TokenValue() string { return b.Token.Value }
func (b Boolean) exprNode()          {}
func (b Boolean) String() string     { return b.Token.Value }

// ============================================================================

type Integer struct {
	Token *token.Token
	Value int64
}

func (i Integer) TokenValue() string { return i.Token.Value }
func (i Integer) exprNode()          {}
func (i Integer) String() string     { return i.Token.Value }

// ============================================================================

type String struct {
	Token *token.Token
	Value string
}

func (s String) TokenValue() string { return s.Token.Value }
func (s String) exprNode()          {}
func (s String) String() string     { return s.Token.Value }

// ============================================================================

type Array struct {
	Token    *token.Token
	Elements []Expr
}

func (a Array) TokenValue() string { return a.Token.Value }
func (a Array) exprNode()          {}
func (a Array) String() string {
	elements := make([]string, 0, len(a.Elements))
	for _, element := range a.Elements {
		elements = append(elements, element.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

// ============================================================================

//{<表达式> : <表达式>, <表达式> : <表达式>, ... }

type Map struct {
	Token    *token.Token
	Elements map[Expr]Expr
}

func (m Map) TokenValue() string { return m.Token.Value }
func (m Map) exprNode()          {}
func (m Map) String() string {
	elements := make([]string, 0, len(m.Elements))
	for key, value := range m.Elements {
		elements = append(elements, key.String()+":"+value.String())
	}
	return "{" + strings.Join(elements, ", ") + "}"
}

// ============================================================================

type Macro struct {
	Token  *token.Token
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStmt
}

func (m Macro) TokenValue() string { return m.Token.Value }
func (m Macro) stmtNode()          {}
func (m Macro) exprNode()          {}
func (m Macro) String() string {
	params := make([]string, 0, len(m.Params))
	for _, param := range m.Params {
		params = append(params, param.String())
	}

	return m.TokenValue() + "" + m.Name.String() + "(" + strings.Join(params, ", ") + ") " + m.Body.String()
}

// ============================================================================

// <前缀运算符><表达式>

type PrefixExpr struct {
	Token    *token.Token
	Operator string
	Right    Expr
}

func (p PrefixExpr) TokenValue() string { return p.Token.Value }
func (p PrefixExpr) exprNode()          {}
func (p PrefixExpr) String() string     { return "(" + p.Operator + p.Right.String() + ")" }

// ============================================================================

// <表达式> <中缀运算符> <表达式>

type InfixExpr struct {
	Token    *token.Token
	Left     Expr
	Operator string
	Right    Expr
}

func (i InfixExpr) TokenValue() string { return i.Token.Value }
func (i InfixExpr) exprNode()          {}
func (i InfixExpr) String() string {
	return "(" + i.Left.String() + " " + i.Operator + " " + i.Right.String() + ")"
}

// ============================================================================

//if (<条件>) <结果> else <可替代的结果>

type IfExpr struct {
	Token       *token.Token
	Condition   Expr
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (i IfExpr) TokenValue() string { return i.Token.Value }
func (i IfExpr) exprNode()          {}
func (i IfExpr) String() string {
	var b strings.Builder
	b.WriteString("if" + i.Condition.String() + " " + i.Consequence.String())
	if i.Alternative != nil {
		b.WriteString("else " + i.Alternative.String())
	}
	return b.String()
}

// ============================================================================

// func <参数列表> <块语句>

type FuncExpr struct {
	Token  *token.Token
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStmt
}

func (f FuncExpr) TokenValue() string { return f.Token.Value }
func (f FuncExpr) exprNode()          {}
func (f FuncExpr) String() string {
	params := make([]string, 0, len(f.Params))
	for _, param := range f.Params {
		params = append(params, param.String())
	}

	return f.TokenValue() + " " + f.Name.String() + " " + "(" + strings.Join(params, ", ") + ") " + f.Body.String()
}

// ============================================================================

// <表达式>(<以逗号分隔的表达式列表>)

type CallExpr struct {
	Token *token.Token
	Func  Expr
	Args  []Expr
}

func (c CallExpr) TokenValue() string { return c.Token.Value }
func (c CallExpr) exprNode()          {}
func (c CallExpr) String() string {
	args := make([]string, 0, len(c.Args))
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}

	return c.Func.TokenValue() + "(" + strings.Join(args, ", ") + ")"
}

// ============================================================================

// <表达式>[<表达式>]

type IndexExpr struct {
	Token *token.Token
	Left  Expr
	Index Expr
}

func (i IndexExpr) TokenValue() string { return i.Token.Value }
func (i IndexExpr) exprNode()          {}
func (i IndexExpr) String() string {
	return "(" + i.Left.String() + "[" + i.Index.String() + "])"
}
