package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/wreulicke/monkey/compiler"
	"github.com/wreulicke/monkey/lexer"
	"github.com/wreulicke/monkey/parser"
	"github.com/wreulicke/monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(bytes.NewBufferString(line))
		p := parser.New(l)

		program := p.Parse()
		fmt.Println(program)
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Exceuting bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		if lastPopped == nil {
			fmt.Fprintln(out, "nil")
			continue
		}
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		fmt.Fprintln(out, "\t", msg)
	}
}
