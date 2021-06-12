package repl

import (
	"bytes"
	"fmt"
	"os"

	prompt "github.com/c-bata/go-prompt"
	"github.com/wreulicke/monkey/lexer"
	"github.com/wreulicke/monkey/token"
)

func Start() {
	p := prompt.New(func(str string) {
		switch str {
		case "exit":
			os.Exit(0)
		default:
			l := lexer.New(bytes.NewBufferString(str))
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		}
	}, func(in prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionPrefix(">> "))
	p.Run()
}
