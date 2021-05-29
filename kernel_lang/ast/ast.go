package ast

import (
	"fmt"
	"strings"
)

type Defs map[string]Tree

type Tree interface{
	String() string
}

// default to wrapping up in struct in case they need to hold line nos

type Int64 struct {
	Val int64
}

type Iden struct {
	Val string
}

type Apply struct {
	Rator Tree
	Rand  Tree
}

type Lambda struct {
	Var  string
	Body Tree
}



// --- printing routines ---
func (d Defs) String() string {
	b := new(strings.Builder)
	for k, v := range d {
		fmt.Fprintf(b, "def %v ", k)
		fmt.Fprintf(b, "%v\n", v.String())
	}
	return b.String()
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

