package repl

import (
	"bufio"
	"fmt"
	"github.com/songzhibin97/mini-interpreter/eval"
	"github.com/songzhibin97/mini-interpreter/object"
	"io"

	"github.com/songzhibin97/mini-interpreter/parser"

	"github.com/songzhibin97/mini-interpreter/lexer"
)

const PROMPT = ">>>"

func Start(in io.Reader, out io.Writer) {
	fmt.Println("Welcome to Mini-interpreter")
	scanner := bufio.NewScanner(in)
	env := object.NewEnv(nil)
	macroEnv := object.NewEnv(nil)
	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		//_, _ = io.WriteString(out, "input >>> "+line+"\r\n")
		p := parser.NewParser(lexer.NewLexer(line))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, s := range p.Errors() {
				_, _ = io.WriteString(out, "\t"+s+"\r\n")
			}
		}
		eval.DefinedMacro(program, macroEnv)
		e := eval.Eval(eval.ExpandMacro(program, macroEnv), env)
		if e == nil || e.Type() == object.NIL {
			continue
		}
		//_, _ = io.WriteString(out, e.Inspect()+"\r\n")
	}
}
