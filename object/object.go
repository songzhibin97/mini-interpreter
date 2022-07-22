package object

import (
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/songzhibin97/mini-interpreter/ast"
)

type Type string

func (t Type) String() string {
	return string(t)
}

type BuiltinFunc func(args ...Object) Object

const (
	INT      Type = "INT"
	String   Type = "STRING"
	BOOL     Type = "BOOL"
	NIL      Type = "NIL"
	RETURN   Type = "RETURN"
	ERROR    Type = "ERROR"
	FUNCTION Type = "FUNCTION"
	BUILTIN  Type = "BUILTIN"
	ARRAY    Type = "ARRAY"
	MAP      Type = "MAP"

	QUOTE Type = "QUOTE"
	MACRO Type = "MACRO"
)

type Object interface {
	Type() Type
	Inspect() string
}

type HashAble interface {
	MapKey() MapKey
}

type MapKey struct {
	Type  Type
	Value uint64
}

// ============================================================================

type Integer struct{ Value int64 }

func (i *Integer) Type() Type      { return INT }
func (i *Integer) Inspect() string { return strconv.Itoa(int(i.Value)) }
func (i *Integer) MapKey() MapKey  { return MapKey{Type: i.Type(), Value: uint64(i.Value)} }

type Stringer struct{ Value string }

func (s Stringer) Type() Type      { return String }
func (s Stringer) Inspect() string { return s.Value }
func (s Stringer) MapKey() MapKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return MapKey{Type: s.Type(), Value: h.Sum64()}
}

type Boolean struct{ Value bool }

func (b *Boolean) Type() Type      { return BOOL }
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *Boolean) MapKey() MapKey {
	if b.Value {
		return MapKey{Type: b.Type(), Value: 1}
	}
	return MapKey{Type: b.Type(), Value: 0}
}

type Nil struct{}

func (n *Nil) Type() Type      { return NIL }
func (n *Nil) Inspect() string { return "nil" }

type Return struct{ Value Object }

func (r *Return) Type() Type      { return RETURN }
func (r *Return) Inspect() string { return r.Value.Inspect() }

type Error struct{ Error string }

func (e *Error) Type() Type      { return ERROR }
func (e *Error) Inspect() string { return e.Error }

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStmt
	Env        *Env
}

func (f *Function) Type() Type { return FUNCTION }
func (f *Function) Inspect() string {
	var b strings.Builder

	params := make([]string, 0, len(f.Parameters))
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	b.WriteString("func " + f.Name.String() + " (" + strings.Join(params, ", ") + ") {\n" + f.Body.String() + "\n}")
	return b.String()
}

type Builtin struct {
	Fn BuiltinFunc
}

func (b *Builtin) Type() Type      { return BUILTIN }
func (b *Builtin) Inspect() string { return "builtin" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type { return ARRAY }
func (a *Array) Inspect() string {
	elements := make([]string, 0, len(a.Elements))
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

type HashValue struct {
	Key   Object
	Value Object
}

type Map struct {
	Elements map[MapKey]HashValue
}

func (m *Map) Type() Type { return MAP }
func (m *Map) Inspect() string {
	pair := make([]string, 0, len(m.Elements))
	for _, value := range m.Elements {
		pair = append(pair, value.Key.Inspect()+":"+value.Value.Inspect())
	}
	return "{" + strings.Join(pair, ", ") + "}"
}

// ============================================================================

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() Type      { return QUOTE }
func (q *Quote) Inspect() string { return "QUOTE(" + q.Node.String() + ")" }

// ============================================================================

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStmt
	Env        *Env
}

func (m *Macro) Type() Type { return MACRO }
func (m *Macro) Inspect() string {
	var b strings.Builder

	params := make([]string, 0, len(m.Parameters))
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	b.WriteString("macro " + " (" + strings.Join(params, ", ") + ") {\n" + m.Body.String() + "\n}")
	return b.String()
}
