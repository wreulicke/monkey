package main

import (
	"log"

	"github.com/wreulicke/monkey/cli"
)

func main() {
	c := cli.New()
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
