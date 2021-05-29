package ast

import (
	"fmt"
	"strings"
)

type Ast interface{
	String() string
}

type Def struct {
	Name string
	Expr Ast
	Next Ast
}
type Nil struct {}

// default to wrapping up in struct in case they need to hold line nos

type Int64 struct {
	Val int64
}

type Iden struct {
	Val string
}

type Apply struct {
	Rator Ast
	Rand  Ast
}

type Lambda struct {
	Var  string
	Body Ast
}

type Cond struct {
	Tests  []Ast
	Consqs []Ast
	Alt    Ast
}



// --- printing routines ---
func (d Def) String() string { 
	b := new(strings.Builder)
	fmt.Fprintf(b, "def %v\n", d.Name)
	fmt.Fprintf(b, "%v\n\n",  d.Expr)
	fmt.Fprintf(b, "%v\n",    d.Next)
	return b.String()
}

func (n Nil) String() string {
	return "end-of-defs"
}

func (i Int64) String() string {
	return fmt.Sprintf("%v", i.Val)
}
func (i Iden) String() string {
	return fmt.Sprintf("%v", i.Val)
}
func (a Apply) String() string {
	return fmt.Sprintf("(@ %v %v)", a.Rator, a.Rand)
}
func (l Lambda) String() string {
	return fmt.Sprintf("(fn %v | %v)", l.Var, l.Body)
}
func (c Cond) String() string {
	b := new(strings.Builder)
	for i, t := range c.Tests {
		fmt.Fprintf(b, "\ncond %v", t)
		fmt.Fprintf(b, "\nthen %v", c.Consqs[i])
	}
	fmt.Fprintf(b, "\nelse %v", c.Alt)
	return b.String()
}


