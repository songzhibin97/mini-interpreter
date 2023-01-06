package ast

type ModifyFunc func(node Node) Node
type Modify func(Node, ModifyFunc) Node

// default

func DefaultModify(node Node, fn ModifyFunc) Node {
	switch n := node.(type) {
	case *Program:
		for i, statement := range n.Stmts {
			n.Stmts[i], _ = DefaultModify(statement, fn).(Stmt)
		}

	case *ExprStmt:
		n.Expr, _ = DefaultModify(n.Expr, fn).(Expr)

	case *ReturnStmt:
		n.Value, _ = DefaultModify(n.Value, fn).(Expr)

	case *VarStmt:
		n.Value, _ = DefaultModify(n.Value, fn).(Expr)

	case *BlockStmt:
		for i, statement := range n.Stmts {
			n.Stmts[i], _ = DefaultModify(statement, fn).(Stmt)
		}

	case *InfixExpr:
		n.Left, _ = DefaultModify(n.Left, fn).(Expr)
		n.Right, _ = DefaultModify(n.Right, fn).(Expr)

	case *PrefixExpr:
		n.Right, _ = DefaultModify(n.Right, fn).(Expr)

	case *IndexExpr:
		n.Left, _ = DefaultModify(n.Left, fn).(Expr)
		n.Index, _ = DefaultModify(n.Index, fn).(Expr)

	case *IfExpr:
		n.Condition, _ = DefaultModify(n.Condition, fn).(Expr)
		n.Consequence, _ = DefaultModify(n.Consequence, fn).(*BlockStmt)
		if n.Alternative != nil {
			n.Alternative, _ = DefaultModify(n.Alternative, fn).(*BlockStmt)
		}

	case *FuncExpr:
		for i, param := range n.Params {
			n.Params[i], _ = DefaultModify(param, fn).(*Identifier)
		}
		n.Body, _ = DefaultModify(n.Body, fn).(*BlockStmt)

	case *Array:
		for i, element := range n.Elements {
			n.Elements[i], _ = DefaultModify(element, fn).(Expr)
		}

	case *Map:
		newElement := make(map[Expr]Expr)
		for k, v := range n.Elements {
			nk, _ := DefaultModify(k, fn).(Expr)
			nv, _ := DefaultModify(v, fn).(Expr)
			newElement[nk] = nv
		}
		n.Elements = newElement
	}

	return fn(node)
}
