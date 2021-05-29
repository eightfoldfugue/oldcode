package main

import (
	"fmt"
    "lang/ast"
	"lang/code"
//	"lang/vm"
)

func main() {

	fmt.Println("here goes")

	tree := ast.ParseFile("program.txt")

	graph := code.Compile(tree)
	fmt.Println(graph)
}
