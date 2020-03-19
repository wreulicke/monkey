package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wreulicke/monkey/ast"
	"github.com/wreulicke/monkey/lexer"
	"github.com/wreulicke/monkey/token"
)

func TestParsingFunctionLiteralWithArrayPattern(t *testing.T) {
	infixTests := []struct {
		input           string
		expectedPattern []ast.Pattern
	}{
		{"fn ([x]) { x }", []ast.Pattern{
			&ast.ArrayPattern{
				Token: token.Token{
					Type:    token.LBRACKET,
					Literal: "[",
				},
				Pattern: []ast.Pattern{
					&ast.Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "x",
						},
						Value: "x",
					},
				},
			}}},
		{"fn ([x, y]) { x }",
			[]ast.Pattern{
				&ast.ArrayPattern{
					Token: token.Token{
						Type:    token.LBRACKET,
						Literal: "[",
					},
					Pattern: []ast.Pattern{
						&ast.Identifier{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "x",
							},
							Value: "x",
						},
						&ast.Identifier{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "y",
							},
							Value: "y",
						},
					},
				}}},
		{"fn ([[x], y]) { x }",
			[]ast.Pattern{
				&ast.ArrayPattern{
					Token: token.Token{
						Type:    token.LBRACKET,
						Literal: "[",
					},
					Pattern: []ast.Pattern{
						&ast.ArrayPattern{
							Token: token.Token{
								Type:    token.LBRACKET,
								Literal: "[",
							},
							Pattern: []ast.Pattern{
								&ast.Identifier{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "x",
									},
									Value: "x",
								},
							},
						},
						&ast.Identifier{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "y",
							},
							Value: "y",
						},
					},
				}}},
		{"fn ({x, y}) { x }",
			[]ast.Pattern{
				&ast.HashPattern{
					Token: token.Token{
						Type:    token.LBRACKET,
						Literal: "{",
					},
					Pattern: []*ast.Identifier{
						&ast.Identifier{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "x",
							},
							Value: "x",
						},
						&ast.Identifier{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "y",
							},
							Value: "y",
						},
					},
				}}},
	}

	for i, tt := range infixTests {
		tt := tt
		t.Run(fmt.Sprintf("tests[%d]", i), func(t *testing.T) {
			l := lexer.New(bytes.NewBufferString(tt.input))
			p := New(l)
			program := p.Parse()
			checkParserErrors(t, p)
			if len(program.Statements) != 1 {
				t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.FunctionLiteral)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
			}
			if len(exp.Parameters) != len(tt.expectedPattern) {
				t.Fatalf("exp.Parameters has not expected length. got=%d, expected:%d", len(exp.Parameters), len(tt.expectedPattern))
			}
			for i, v := range tt.expectedPattern {
				if exp.Parameters[i].TokenLiteral() != v.TokenLiteral() {
					t.Fatalf("exp.Parameters[i].TokenLiteral is not expected. got=%s, want=%s", exp.Parameters[i].TokenLiteral(), v.TokenLiteral())
				}
				if exp.Parameters[i].String() != v.String() {
					t.Fatalf("exp.Parameters[i].String is not expected. got=%s, want=%s", exp.Parameters[i].String(), v.String())
				}
			}
		})
	}

}

func TestParsingPipelineOperator(t *testing.T) {
	infixTests := []struct {
		input      string
		leftString string
	}{
		{"5 | fn (x) { x }", "5"},
		{"(5 + 10) | fn (x) { x }", "(5 + 10)"},
		{"5 + 10 | fn (x) { x }", "(5 + 10)"},
		{"x(10 + 2) - 3 | fn (x) { x }", "(x((10 + 2)) - 3)"},
		{"x(10 + 2) - (3 + 2) | fn (x) { x }", "(x((10 + 2)) - (3 + 2))"},
		{"x(10 + 2) - (3 | fn (x) { x } ) | fn (x) { x }", "(x((10 + 2)) - (3 | fn(x) x))"},
	}

	for _, tt := range infixTests {
		l := lexer.New(bytes.NewBufferString(tt.input))
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != "|" {
			t.Fatalf("exp.Operator is not '%s'. got=%s", "|", exp.Operator)
		}
		if exp.Left.String() != tt.leftString {
			t.Errorf("exp.String() is not %q. got=%q", tt.leftString, exp.Left.String())
		}
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := `{}`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", k)
		}
		expectedValue := expected[literal.String()]
		testNumberLiteral(t, v, fmt.Sprintf("%v", expectedValue))
	}
}

func TestParsingHashLiteralWithExpression(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]func(e ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "0",
				},
				Value: "0",
			}, "+", &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "1",
				},
				Value: "1",
			})
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "10",
				},
				Value: "10",
			}, "-", &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "8",
				},
				Value: "8",
			})
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "15",
				},
				Value: "15",
			}, "/", &ast.NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "5",
				},
				Value: "5",
			})
		},
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", k)
		}
		assertion, ok := expected[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		assertion(v)
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IndexExpression. got=%T", stmt.Expression)
	}
	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "1",
		},
		Value: "1",
	}, "+", &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "1",
		},
		Value: "1",
	})
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	testNumberLiteral(t, literal.Elements[0], "1")
	testInfixExpression(t, literal.Elements[1], &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "2",
		},
		Value: "2",
	}, "*", &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "2",
		},
		Value: "2",
	})
	testInfixExpression(t, literal.Elements[2], &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "3",
		},
		Value: "3",
	}, "+", &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "3",
		},
		Value: "3",
	})
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world"`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Fatalf("function arguments wrong. want 3, got =%d", len(call.Arguments))
	}
	testExpression(t, call.Arguments[0], &ast.NumberLiteral{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "1",
		},
		Value: "1",
	})
	testInfixExpression(t, call.Arguments[1],
		&ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "2",
			},
			Value: "2",
		},
		"*",
		&ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "3",
			},
			Value: "3",
		})
	testInfixExpression(t, call.Arguments[2],
		&ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "4",
			},
			Value: "4",
		},
		"+",
		&ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "5",
			},
			Value: "5",
		})
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got =%d", len(fn.Parameters))
	}
	testExpression(t, fn.Parameters[0].(*ast.Identifier), &ast.Identifier{
		Token: token.Token{
			Type:    token.IDENT,
			Literal: "x",
		},
		Value: "x",
	})
	testExpression(t, fn.Parameters[1].(*ast.Identifier), &ast.Identifier{
		Token: token.Token{
			Type:    token.IDENT,
			Literal: "y",
		},
		Value: "y",
	})
	if len(fn.Body.Statements) != 1 {
		t.Fatalf("fn.Body has not 1 statements. got=%d", len(fn.Body.Statements))
	}
	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression,
		&ast.Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "x",
			},
			Value: "x",
		},
		"+",
		&ast.Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "y",
			},
			Value: "y",
		})
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition,
		&ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "x",
		}, Value: "x"},
		"<", &ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "y",
		}, Value: "y"}) {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative is not nil. got=%T", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition,
		&ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "x",
		}, Value: "x"},
		"<", &ast.Identifier{Token: token.Token{
			Type:    token.IDENT,
			Literal: "y",
		}, Value: "y"}) {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence does not have 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative == nil {
		t.Errorf("exp.Alternative is nil.")
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("alternative does not have 1 statements. got=%d", len(exp.Consequence.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 > 4 != 3 < 4",
			"((5 > 4) != (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(bytes.NewBufferString(tt.input))
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  string
		operator   string
		rightValue string
	}{
		{"5 + 5", "5", "+", "5"},
		{"5 - 5", "5", "-", "5"},
		{"5 * 5", "5", "*", "5"},
		{"5 / 5", "5", "/", "5"},
		{"5 > 5", "5", ">", "5"},
		{"5 < 5", "5", "<", "5"},
		{"5 == 5", "5", "==", "5"},
		{"5 != 5", "5", "!=", "5"},
	}

	for _, tt := range infixTests {
		l := lexer.New(bytes.NewBufferString(tt.input))
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testNumberLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if !testNumberLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    string
	}{
		{"!5", "!", "5"},
		{"-5", "-", "5"},
	}
	for _, tt := range prefixTests {
		l := lexer.New(bytes.NewBufferString(tt.input))
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testNumberLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func testNumberLiteral(t *testing.T, il ast.Expression, value string) bool {
	l, ok := il.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("il not ast.NumberLiteral. got=%T", il)
		return false
	}
	if l.Value != value {
		t.Errorf("l.Value not %s. got=%s", l.Value, value)
		return false
	}
	return true
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Expression)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "foobar", literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	testIdentifier(t, stmt.Expression, "foobar")
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}

}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      ast.Expression
	}{
		{"let x = 5;", "x", &ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "5",
			},
			Value: "5",
		}},

		{"let y = true;", "y", &ast.BooleanLiteral{
			Token: token.Token{
				Type:    token.TRUE,
				Literal: "true",
			},
			Value: true,
		}},

		{"let foobar = y;", "foobar", &ast.Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "y",
			},
			Value: "y",
		}},
	}

	for _, tt := range tests {
		l := lexer.New(bytes.NewBufferString(tt.input))
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			continue
		}

		val := stmt.(*ast.LetStatement).Value
		if !testExpression(t, val, tt.expectedValue) {
			continue
		}
	}
}

func TestLetStatementWithArrayPattern(t *testing.T) {
	input := `let [x, y] = [1, 2];`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("Parse returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have 1 statement. got=%d", len(program.Statements))
	}

	s := program.Statements[0]
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral is not 'let'. got=%q", s.TokenLiteral())
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s is not *ast.LetStatement")
	}

	patterns, ok := letStmt.Pattern.(*ast.ArrayPattern)
	if !ok {
		t.Errorf("letStmt.Pattern is not ArrayPattern. got=%T", letStmt.Pattern)
	}

	testIdentifier(t, patterns.Pattern[0].(*ast.Identifier), "x")
	testIdentifier(t, patterns.Pattern[1].(*ast.Identifier), "y")
}

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(bytes.NewBufferString(input))
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("Parse returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testExpression(t *testing.T, exp ast.Expression, expected ast.Expression) bool {
	switch v := expected.(type) {
	case *ast.NumberLiteral:
		return testNumberLiteral(t, exp, v.Value)
	case *ast.BooleanLiteral:
		return testBooleanLiteral(t, exp, v.Value)
	case *ast.Identifier:
		return testIdentifier(t, exp, v.Value)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left ast.Expression, operator string, right ast.Expression) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", exp)
		return false
	}
	if infixExp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s", operator, infixExp.Operator)
	}
	if !testExpression(t, infixExp.Left, left) {
		return false
	}
	if !testExpression(t, infixExp.Right, right) {
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}
	if b.Value != value {
		t.Errorf("ident.Value not %t. got=%t", value, b.Value)
		return false
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("ident.TokenLiteral not %t. got=%s", value, b.TokenLiteral())
		return false
	}
	return true

}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement")
	}

	if letStmt.Pattern.(*ast.Identifier).Value != name {
		t.Errorf("letStmt.Name.Value is not '%s'. got=%s", name, letStmt.Pattern.(*ast.Identifier).Value)
		return false
	}

	if letStmt.Pattern.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s'. got=%s", name, letStmt.Pattern.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %+v", err)
	}
	t.FailNow()
}
