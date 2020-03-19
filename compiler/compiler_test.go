package compiler

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wreulicke/monkey/lexer"
	"github.com/wreulicke/monkey/parser"
)

func TestCompile(t *testing.T) {
	input := `
	let y = "Hello World!! tes\n" 
	let x = "test\n";
	printf(y)
	printf(x)
	`
	l := lexer.New(bytes.NewBufferString(input))
	p := parser.New(l)
	program := p.Parse()
	m, err := Compile(program)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(m)
	t.Error("teest")

}
