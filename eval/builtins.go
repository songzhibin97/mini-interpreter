package eval

import (
	"fmt"

	"github.com/songzhibin97/mini-interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return &object.Error{Error: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
		}
		switch arg := args[0].(type) {
		case *object.Map:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.Stringer:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return &object.Error{Error: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
		}
	}},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return &object.Nil{}
		},
	},
}
