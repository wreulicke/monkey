package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/wreulicke/monkey/ast"
	"github.com/wreulicke/monkey/code"
)

var typeNames = []string{
	"INTEGER",
	"BOOLEAN",
	"STRING",
	"ARRAY",
	"HASH",
	"FUNCTION",
	"NULL",
	"RETURN",
	"ERROR",
	"BUILTIN",
}

type ObjectType int

const (
	INTEGER ObjectType = iota
	BOOLEAN
	STRING
	ARRAY
	HASH
	FUNCTION
	NULL
	RETURN
	ERROR
	BUILTIN
	COMPILED_FUNCTION
	CLOSURE
)

func (o ObjectType) String() string {
	return typeNames[o]
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (n *Integer) Type() ObjectType {
	return INTEGER
}

func (n *Integer) Inspect() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *Integer) HashKey() HashKey {
	return HashKey{Type: n.Type(), Value: uint64(n.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, v := range a.Elements {
		elements = append(elements, v.Inspect())
	}
	out.WriteRune('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteRune(']')
	return out.String()
}

type Null struct {
}

func (b *Null) Type() ObjectType {
	return NULL
}

func (b *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []ast.Pattern
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, v := range f.Parameters {
		params = append(params, v.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (f *Builtin) Type() ObjectType {
	return BUILTIN
}

func (f *Builtin) Inspect() string {
	return "builtin function"
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType { return CLOSURE }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}
