package ast

type ModifyFunc func(node Node) Node

func Modify(node Node, fn ModifyFunc) Node {
	switch n := node.(type) {
	case *Program:
		for i, statement := range n.Stmts {
			n.Stmts[i], _ = Modify(statement, fn).(Stmt)
		}

	case *ExprStmt:
		n.Expr, _ = Modify(n.Expr, fn).(Expr)

	case *ReturnStmt:
		n.Value, _ = Modify(n.Value, fn).(Expr)

	case *VarStmt:
		n.Value, _ = Modify(n.Value, fn).(Expr)

	case *BlockStmt:
		for i, statement := range n.Stmts {
			n.Stmts[i], _ = Modify(statement, fn).(Stmt)
		}

	case *InfixExpr:
		n.Left, _ = Modify(n.Left, fn).(Expr)
		n.Right, _ = Modify(n.Right, fn).(Expr)

	case *PrefixExpr:
		n.Right, _ = Modify(n.Right, fn).(Expr)

	case *IndexExpr:
		n.Left, _ = Modify(n.Left, fn).(Expr)
		n.Index, _ = Modify(n.Index, fn).(Expr)

	case *IfExpr:
		n.Condition, _ = Modify(n.Condition, fn).(Expr)
		n.Consequence, _ = Modify(n.Consequence, fn).(*BlockStmt)
		if n.Alternative != nil {
			n.Alternative, _ = Modify(n.Alternative, fn).(*BlockStmt)
		}

	case *FuncExpr:
		for i, param := range n.Params {
			n.Params[i], _ = Modify(param, fn).(*Identifier)
		}
		n.Body, _ = Modify(n.Body, fn).(*BlockStmt)

	case *Array:
		for i, element := range n.Elements {
			n.Elements[i], _ = Modify(element, fn).(Expr)
		}

	case *Map:
		newElement := make(map[Expr]Expr)
		for k, v := range n.Elements {
			nk, _ := Modify(k, fn).(Expr)
			nv, _ := Modify(v, fn).(Expr)
			newElement[nk] = nv
		}
		n.Elements = newElement
	}

	return fn(node)
}
