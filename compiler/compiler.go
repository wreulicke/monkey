package compiler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/wreulicke/monkey/ast"
)

type Environment struct {
	env map[string]value.Value
}

func (e *Environment) Get(key string) (value.Value, bool) {
	v, ok := e.env[key]
	return v, ok
}

func (e *Environment) Put(key string, v value.Value) {
	e.env[key] = v
}

func Compile(program *ast.Program) (*ir.Module, error) {
	m := ir.NewModule()
	x := m.NewFunc("main", types.Void)
	block := x.NewBlock("")
	printf := m.NewFunc("printf", types.Void)
	printf.Sig.Variadic = true
	env := &Environment{
		env: map[string]value.Value{},
	}
	env.Put("printf", printf)
	for _, s := range program.Statements {
		inst, err := translateInsts(m, block, env, s)
		if err != nil {
			return nil, err
		} else if inst != nil {
			block.Insts = append(block.Insts, inst)
		}
	}

	block.NewRet(nil)
	return m, nil
}

func translateInsts(m *ir.Module, block *ir.Block, env *Environment, node ast.Node) (ir.Instruction, error) {
	switch node := node.(type) {
	case *ast.LetStatement:
		i := node.Pattern.(*ast.Identifier)
		v, err := translateValue(m, block, node.Value, env)
		if err != nil {
			return nil, err
		}
		aloc := block.NewAlloca(v.Type())
		env.Put(i.Value, aloc)
		block.NewStore(v, aloc)
		return nil, nil
	case *ast.ExpressionStatement:
		return translateInsts(m, block, env, node.Expression)
	case *ast.CallExpression:
		args := []value.Value{}
		for _, a := range node.Arguments {
			v, err := translateValue(m, block, a, env)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		i, ok := node.Function.(*ast.Identifier)
		if !ok {
			return nil, errors.New("not implemented")
		}

		f, ok := env.Get(i.Value)
		if !ok {
			return nil, errors.New("not implemented")
		}
		return ir.NewCall(f, args...), nil
	}

	return nil, errors.New("not implemented")
}

func translateValue(m *ir.Module, block *ir.Block, node ast.Node, env *Environment) (value.Value, error) {
	switch node := node.(type) {
	case *ast.ArrayLiteral:
		panic("not supported")
	case *ast.NumberLiteral:
		i, err := strconv.ParseInt(node.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return constant.NewInt(types.I64, i), nil
	case *ast.StringLiteral:
		ch := constant.NewCharArrayFromString(node.Value)
		return m.NewGlobalDef("$const_"+node.String(), ch), nil
	case *ast.FunctionLiteral:
	case *ast.Identifier:
		v, ok := env.Get(node.Value)
		if !ok {
			return nil, fmt.Errorf("%s is not found", node.Value)
		}
		return block.NewLoad(v.(*ir.InstAlloca).ElemType, v), nil
	}
	panic("not implemented")
}
