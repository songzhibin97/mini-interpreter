package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expr { return &Integer{Value: 1} }
	two := func() Expr { return &Integer{Value: 2} }
	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*Integer)
		if !ok {
			return node
		}
		if integer.Value != 1 {
			return node
		}
		integer.Value = 2
		return integer
	}
	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Stmts: []Stmt{&ExprStmt{Expr: one()}}},
			&Program{
				Stmts: []Stmt{&ExprStmt{Expr: two()}},
			},
		},
		{
			&InfixExpr{Left: one(), Operator: "+", Right: two()},
			&InfixExpr{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpr{Left: two(), Operator: "+", Right: one()},
			&InfixExpr{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpr{Operator: "-", Right: one()},
			&PrefixExpr{Operator: "-", Right: two()},
		},
		{
			&IndexExpr{Left: one(), Index: one()},
			&IndexExpr{Left: two(), Index: two()},
		},
		{
			&IfExpr{
				Condition: one(),
				Consequence: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: one()},
					}},
				Alternative: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: one()},
					},
				}},
			&IfExpr{
				Condition: two(),
				Consequence: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: two()},
					}},
				Alternative: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: two()},
					},
				},
			},
		},
		{
			&ReturnStmt{Value: one()},
			&ReturnStmt{Value: two()},
		},
		{
			&VarStmt{Value: one()},
			&VarStmt{Value: two()},
		},
		{
			&FuncExpr{
				Params: []*Identifier{},
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: one()},
					}},
			},
			&FuncExpr{
				Params: []*Identifier{},
				Body: &BlockStmt{
					Stmts: []Stmt{
						&ExprStmt{Expr: two()},
					}},
			},
		},
		{
			&Array{Elements: []Expr{one(), one()}},
			&Array{Elements: []Expr{two(), two()}},
		},
	}
	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)
		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}
	mp := &Map{
		Elements: map[Expr]Expr{
			one(): one(),
			one(): one(),
		},
	}
	Modify(mp, turnOneIntoTwo)
	for key, val := range mp.Elements {
		key, _ := key.(*Integer)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}
		val, _ := val.(*Integer)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
