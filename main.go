package main

import (
	"github.com/motoki317/lisp-interpreter/lisp"
	"github.com/motoki317/lisp-interpreter/node"
	"github.com/motoki317/lisp-interpreter/token"
	"os"
)

func main() {
	i := lisp.NewInterpreter(node.NewParser(token.NewTokenizer(os.Stdin)), os.Stdout)
	i.ReadLoop()
}
