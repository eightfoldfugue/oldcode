package code

import (
	"fmt"
	"strings"
)


type Code struct {
	SizeGlobals int
	Ops         []Op
}

type Op int64

const (
	Halt Op = iota

	Add
	Sub
	Mul

	GetDef
	SetDef
	
	PushInt
	PushFun

	PushMark
	Grab
	Access
	Call
	Return
)


type graph map[string] *block

type block struct {
	code asm
	refs refSet
	addr int
}

type refSet map[string] []int

type asm []word

type word struct {
	name string
	datum Op
}

// additional names
const imm_val = "imm_val"
const label = "label"

func (o Op) toString() string {
	return [...]string{
		"Halt",

		"Add",
		"Sub",
		"Mul",

		"GetDef",
		"SetDef",

		"PushInt",
		"PushFun",

		"PushMark",
		"Grab",
		"Access",
		"Call",
		"Return",
	}[o]
}

// block operations
func newblock() *block {
	return &block{asm{}, refSet{}, 0}
}

// operations on asm part
func (b *block) emit(o Op) {
	b.code = append(b.code, word{o.toString(), o})
}
func (b *block) immeadiate(o Op) {
	b.code = append(b.code, word{imm_val, o})
}

func (b *block) label() int {
	b.code = append(b.code, word{"label", Op(0)})
	return len(b.code) -1
}
func (b *block) backpatch(label_index, set_value int) {
	b.code[label_index].datum = Op(set_value)
}
func (b *block) last() int {
	return len(b.code) - 1
}

// operations on refSet part
func (b *block) addRef(name string, loc int) {
	locs, inSet := b.refs[name]
	if inSet {
		locs = append(locs, loc)
		b.refs[name] = locs
	} else {
		b.refs[name] = []int{loc}
	}
}



// printing routines
func (a asm) String() string {
	// header
	b := new(strings.Builder)
	fmt.Fprint(b, "\n")
	header := "%6v | %-11v | %-6v \n"
	fmt.Fprintf(b, header, "index", "instruction", "datum")
	for i := 0; i < 29 ; i++ {
		fmt.Fprint(b, "-")
	}
	fmt.Fprint(b, "\n")
	
	for i, v := range a {
		fmt.Fprintf(b, "%6v |", i)
		fmt.Fprintf(b, " %-11v |", v.name)

		// only print relevant data
		switch v.name {
		case imm_val, label:
			fmt.Fprintf(b, " %v", v.datum)
		}
		fmt.Fprint(b, "\n")
	}
	return b.String()
}

func (b block) String() string {
	s := new(strings.Builder)

	s.WriteString(b.code.String())

	fmt.Fprintf(s, "and references:\n%v\n", b.refs)
	fmt.Fprintf(s, "with address: %v\n", b.addr)

	return s.String()
}

func (g graph) String() string {
	s := new(strings.Builder)
	for k, v := range g {
		fmt.Fprintf(s, "%v has intructions:", k)
		fmt.Fprintf(s, "%v\n\n", v)
	}
	return s.String()
}
