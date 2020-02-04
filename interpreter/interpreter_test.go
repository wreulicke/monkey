package interpreter

import (
	"fmt"
	"testing"

	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/lexer"
	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/object"
	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/parser"
)

func TestPipelineOperatorExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`1 | fn(x) { x }`, 1},
		{`[1, 2] | fn(x) { x[0] + x[1] }`, 3},
		{`[1, 2] | fn(x) { x[0] + x[1] } | fn(x) {x * 2}`, 6},
		// Array Pattern
		{`[1, 2] | fn([x, y]) { x + y }`, 3},
		{`[[1], 2] | fn([[x], y]) { x + y }`, 3},
		// HashPattern
		{`{"x": 1, "y": 2} | fn({x, y}) { x + y }`, 3},

		// builtins
		{`[1, 2] | len`, 2},
		{`[1, 2] | first`, 1},

		{`
		let _first = fn(x) { x[0] };
		[1, 2] | _first
		`, 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		}
	}

}

func TestHashLiterals(t *testing.T) {
	input := `
	let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didnt return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Error("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestMap(t *testing.T) {
	input := `
	let map = fn(arr, f) {
		let iter = fn(arr, accumulated) {
			if(len(arr) == 0) {
				accumulated
			} else {
				iter(rest(arr), push(accumulated, f(first(arr))))
			}
		}
		iter(arr, [])
	}
	let a = [1, 2, 3, 4]
	let double = fn(x) { x * 2}
	map(a, double)
	`
	evaluated := testEval(input)
	arr := evaluated.(*object.Array)
	testIntegerObject(t, arr.Elements[0], 2)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
	testIntegerObject(t, arr.Elements[3], 8)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"let i = 0; [1][i]", 1},
		{"[1, 2, 3][1 + 1]", 3},
		{"let myArray = [1, 2, 3]; myArray[2]", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}

}

func TestFibbo(t *testing.T) {
	input := `
	let fibb = fn(x) { 
		if(x == 0) { 
			x 
		} 
		else {
			if (x == 1) { x }
			else { fibb(x - 1) + fibb(x - 2) }
		}
	};
	fibb(2)
	`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 1)
}

func TestStringEqaulityExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"Hello" == "Hello"`, true},
		{`"Hello" == "World"`, false},
		{`"Hello" + "World" == "World"`, false},
		{`"Hello" + "World" == "HelloWorld"`, true},
		{`"Hello" + " " + "World" == "World"`, false},
		{`"Hello" + " " + "World" == "Hello World"`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello"+ " " + "World"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String was wrong value. got=%q", str.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String was wrong value. got=%q", str.Value)
	}
}

func TestClosure(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) {x + y}
}
let addTwo = newAdder(2)
addTwo(2)
`
	testIntegerObject(t, testEval(input), 4)
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%t (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Errorf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},

		{"let [x, y] = [15, 0]; x", 15},

		// {`
		// 	let x = 15
		// 	let y = 10
		// 	if (x == 15) {
		// 		let y = 0
		// 		if (y == 10) {
		// 			return 15
		// 		}
		// 	}
		// 	return y
		// `, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if(10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return true + false
				}
				return 1
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier is not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }]`,
			"unusable as hash key: FUNCTION",
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("tests[%d]", i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Log(tt.input)
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				return
			}
			if errObj.Message != tt.expectedMessage {
				t.Log(tt.input)
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
			}
		})
	}

}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10
				}
				return 1
			}`,
			10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(true) { 10 }", 10},
		{"if(false) { 10 }", nil},
		{"if (1) { 1 }", 1},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if v, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(v))
		} else {
			testNullObject(t, evaluated)
		}
	}

}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Message != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Message, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj == NULL {
		return true
	}
	t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	return false
}
