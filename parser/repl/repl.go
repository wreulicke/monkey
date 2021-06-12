package repl

import (
	"bytes"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/wreulicke/monkey/lexer"
	"github.com/wreulicke/monkey/parser"
)

func Start() {
	p := prompt.New(func(str string) {
		switch str {
		case "exit":
			os.Exit(0)
		default:
			l := lexer.New(bytes.NewBufferString(str))
			p := parser.New(l)

			program := p.Parse()
			if len(p.Errors()) != 0 {
				printParseError(p.Errors())
				return
			}

			fmt.Println(program.String())
		}
	}, func(in prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionPrefix(">> "))
	p.Run()
}

func printParseError(errors []error) {
	for _, msg := range errors {
		fmt.Println("\t", msg)
	}
}
