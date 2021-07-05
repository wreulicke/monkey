package cli

import (
	"os"

	"github.com/spf13/cobra"

	interpreterRepl "github.com/wreulicke/monkey/interpreter/repl"
	lexerRepl "github.com/wreulicke/monkey/lexer/repl"
	parserRepl "github.com/wreulicke/monkey/parser/repl"
	vmRepl "github.com/wreulicke/monkey/vm/repl"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "monkey",
		Short: "monkey interpreter",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	c.AddCommand(NewInterpreterCommand(), NewLexerCommand(), NewParserCommand(), NewVMCommand())
	return c
}

func NewInterpreterCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "interpreter",
		Short: "interpreter repl",
		Run: func(cmd *cobra.Command, args []string) {
			interpreterRepl.Start()
		},
	}
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

func NewVMCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "parser",
		Short: "parser repl",
		Run: func(cmd *cobra.Command, args []string) {
			parserRepl.Start()
		},
	}
	return c
}

func NewParserCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm",
		Short: "vm repl",
		Run: func(cmd *cobra.Command, args []string) {
			vmRepl.Start(os.Stdin, os.Stdout)
		},
	}
	return c
}
