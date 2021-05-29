package main

import (
	"fmt"
    "lang/ast"
	"lang/tok"
	"lang/code"
	"lang/vm"
)

func main() {
	toks := tok.LexFile("program.txt")
	defs := ast.Parse(toks)
	fmt.Println(defs)


	grph := code.Compile(defs)
	fmt.Print(grph)


	inst := code.Assemble(grph)
	fmt.Println(inst)

	
	answ := vm.Eval(inst)
    fmt.Println(answ)

}
