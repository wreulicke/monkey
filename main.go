package main

import (
	"log"

	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/cli"
)

func main() {
	c := cli.New()
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
