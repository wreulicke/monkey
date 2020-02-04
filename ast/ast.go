package ast

import (
	"bytes"
	"strings"

	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/token"
)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type statement struct{}

func (s *statement) statementNode() {}

type expression struct{}

func (e *expression) expressionNode() {}

type pattern struct{}

func (p *pattern) patternNode() {}

type LetStatement struct {
	statement
	Token   token.Token
	Pattern Pattern
	Value   Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteRune(' ')
	out.WriteString(ls.Pattern.String())

	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	statement
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	out.WriteRune(' ')
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type BlockStatement struct {
	statement
	Token      token.Token
	Statements []Statement
}

func (ie *BlockStatement) TokenLiteral() string {
	return ie.Token.Literal
}

func (p *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	statement
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	expression
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteRune(')')

	return out.String()
}

type InfixExpression struct {
	expression
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (pe *InfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(pe.Left.String())
	out.WriteRune(' ')
	out.WriteString(pe.Operator)
	out.WriteRune(' ')
	out.WriteString(pe.Right.String())
	out.WriteRune(')')

	return out.String()
}

type IfExpression struct {
	expression
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteRune(' ')
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type CallExpression struct {
	expression
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteRune('(')
	out.WriteString(strings.Join(args, ", "))
	out.WriteRune(')')

	return out.String()
}

type IndexExpression struct {
	expression
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(ie.Left.String())
	out.WriteRune('[')
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type Identifier struct {
	expression
	pattern
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	expression
	Token token.Token
	Value string
}

func (i *NumberLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *NumberLiteral) String() string {
	return i.Value
}

type BooleanLiteral struct {
	expression
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func (b *StringLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *StringLiteral) String() string {
	return b.Token.Literal
}

type ArrayLiteral struct {
	expression
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range al.Elements {
		args = append(args, a.String())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(args, ", "))
	out.WriteRune(']')

	return out.String()
}

type HashLiteral struct {
	expression
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}

type FunctionLiteral struct {
	expression
	Token      token.Token
	Parameters []Pattern
	Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteRune('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Pattern interface {
	Node
	patternNode()
}

type ArrayPattern struct {
	pattern
	Token   token.Token
	Pattern []Pattern
}

func (ap *ArrayPattern) TokenLiteral() string {
	return ap.Token.Literal
}

func (ap *ArrayPattern) String() string {
	var out bytes.Buffer

	patterns := []string{}
	for _, a := range ap.Pattern {
		patterns = append(patterns, a.String())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(patterns, ", "))
	out.WriteRune(']')

	return out.String()
}

type HashPattern struct {
	pattern
	Token   token.Token
	Pattern []*Identifier
}

func (hp *HashPattern) TokenLiteral() string {
	return hp.Token.Literal
}

func (hp *HashPattern) String() string {
	var out bytes.Buffer

	patterns := []string{}
	for _, p := range hp.Pattern {
		patterns = append(patterns, p.String())
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(patterns, ", "))
	out.WriteRune('}')

	return out.String()
}
