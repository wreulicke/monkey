package interpreter

import (
	"fmt"
	"strconv"

	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/ast"
	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		bindPattern(env, node.Pattern, val)
		return val
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		return evalCallExpression(fn, node.Arguments, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.ArrayLiteral:
		elements, err := evalExpressions(node.Elements, env)
		if err != nil {
			return err
		}
		return &object.Array{
			Elements: elements,
		}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.NumberLiteral:
		i, err := strconv.ParseInt(node.Value, 10, 64)
		if err != nil {
			return newError("cannot convert int. %s", node.Value)
		}
		return &object.Integer{Value: i}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.FunctionLiteral:
		f := &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		return f
	}
	return newError("Unsupported ast.Node got=%T", node)
}

func bindPattern(env *object.Environment, pattern ast.Pattern, val object.Object) *object.Error {
	switch node := pattern.(type) {
	case *ast.Identifier:
		env.Set(node.Value, val)
	case *ast.ArrayPattern:
		array, ok := val.(*object.Array)
		if !ok {
			return newError("initializer is not Array. cannot destruct. got=%T", val)
		}
		for idx, v := range node.Pattern {
			if idx < len(array.Elements) {
				bindPattern(env, v, array.Elements[idx])
			}
		}
	case *ast.HashPattern:
		hash, ok := val.(*object.Hash)
		if !ok {
			return newError("initializer is not Hash. cannot destruct. got=%T", val)
		}
		for _, v := range node.Pattern {
			string := object.String{Value: v.String()}
			hashKey := string.HashKey()
			pair, ok := hash.Pairs[hashKey]
			if ok {
				bindPattern(env, v, pair.Value)
			}
		}

	}
	return nil
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	var result []object.Object
	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return nil, evaluated
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalCallExpression(fn object.Object, arguments []ast.Expression, env *object.Environment) object.Object {
	args, err := evalExpressions(arguments, env)
	if err != nil {
		return err
	}
	return callFunction(fn, args)
}

func callFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		functionEnv := extendFunctionEnv(fn, args)
		return unwrapReturnValue(Eval(fn.Body, functionEnv))
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return newError("not a function: %s", fn.Type())
}

func extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := function.Env.NewEnclosedEnvironment()

	for paramIdx, param := range function.Parameters {
		fmt.Printf("%T %s\n", param, param.String())
		bindPattern(env, param, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(o object.Object) object.Object {
	switch o := o.(type) {
	case *object.ReturnValue:
		return o.Value
	default:
		return o
	}
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s, env)
		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			return v
		}
	}

	return result

}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s, env)
		if result != nil {
			switch result.Type() {
			case object.RETURN:
				return result
			case object.ERROR:
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func isTruthy(o object.Object) bool {
	switch o {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case right.Type() == object.FUNCTION || right.Type() == object.BUILTIN:
		return evalPipelineOperator(left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPipelineOperator(left object.Object, right object.Object) object.Object {
	return callFunction(right, []object.Object{left})
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		hashKey, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", index.Type())
		}
		return evalHashIndexExpression(left, hashKey.HashKey())
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	leftValue := left.(*object.Array)
	indexValue := index.(*object.Integer)
	if indexValue.Value < 0 || indexValue.Value >= int64(len(leftValue.Elements)) {
		return NULL
	}
	return leftValue.Elements[indexValue.Value]
}

func evalHashIndexExpression(left object.Object, hashKey object.HashKey) object.Object {
	leftValue := left.(*object.Hash)
	r, ok := leftValue.Pairs[hashKey]
	if !ok {
		return NULL
	}
	return r.Value
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return FALSE
	default:
		return FALSE
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := map[object.HashKey]object.HashPair{}

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}
	return &object.Hash{Pairs: pairs}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if v, ok := env.Get(node.Value); ok {
		return v
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier is not found: %s", node.Value)
}

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}
