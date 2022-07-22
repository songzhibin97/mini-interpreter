package main

import (
	"os"

	"github.com/songzhibin97/mini-interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
