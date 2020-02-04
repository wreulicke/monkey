package cli

import (
	"github.com/spf13/cobra"

	interpreterRepl "github.com/wreulicke/monkey/interpreter/repl"
	lexerRepl "github.com/wreulicke/monkey/lexer/repl"
	parserRepl "github.com/wreulicke/monkey/parser/repl"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "monkey",
		Short: "monkey interpreter",
		Run: func(cmd *cobra.Command, args []string) {
			interpreterRepl.Start()
		},
	}
	c.AddCommand(NewLexerCommand(), NewParserCommand())
	return c
}

func NewLexerCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "lexer",
		Short: "lexer repl",
		Run: func(cmd *cobra.Command, args []string) {
			lexerRepl.Start()
		},
	}
	return c
}

func NewParserCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "parser",
		Short: "parser repl",
		Run: func(cmd *cobra.Command, args []string) {
			parserRepl.Start()
		},
	}
	return c
}
