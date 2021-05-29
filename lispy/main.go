package main

import (
	"fmt"
	"log"
)

func main() {
	var n1, n2, n3 word
	n1 = 6
	n2 = 7
	n3 = 8
	n1 = setFlag(n1, fixNum)
	n2 = setFlag(n2, fixNum)
	n3 = setFlag(n3, fixNum)

	vm := new(machine)
	vm.push(n1)
	vm.makeLeaf()
	vm.push(n2)
	vm.makeLeaf()

	vm.makeLink()
	
	vm.push(n3)
	vm.makeLeaf()

	vm.makeLink()

	vm.showStack()
	fmt.Println(vm.links)
}

func show(w word) {
	fmt.Printf("%064b\n", w)
}

type machine struct {
	sp    int
	stack [4]word

	free  word
	scan  word
	root  word
	
	links [4]link

}


type word uint64

const reservedFlags word = 8
const initIota word = 64 - reservedFlags
const maxNum word = fixNum-1
const (
	fixNum word = 1 << (initIota+iota)
	nil
	leaf
	pair
	gcMark
)

type link struct {
	size word
	car  word
	cdr  word
}


func setFlag(b, flag word) word    { return b | flag }
func clearFlag(b, flag word) word  { return b &^ flag }

// toggle: if its on, turn it off, if its off, turn it on
func toggleFlag(b, flag word) word { return b ^ flag }
func hasFlag(b, flag word) bool    { return b&flag != 0 }

func (vm *machine) showStack() {
	for i := 0; i < vm.sp; i++ {
		fmt.Printf("%064b\n", vm.stack[i])
	}
}
func (vm *machine) push(w word) {
	vm.stack[vm.sp] = w
	vm.sp++
}
func (vm *machine) pop() word {
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *machine) makeLeaf() {
	el := vm.pop()
	lf := setFlag(el, leaf)
	vm.push(lf)
}

func (vm *machine) makeLink() {
	l := vm.pop()
	r := vm.pop()
	sizel := vm.sizeOf(l)
	sizer := vm.sizeOf(r)
	if sizel == 0 {
		vm.push(r) // nil is the linking identity val, r is unchanged
	} else if sizer == 0 {
		vm.push(l) // symmetric case
	} else {
		size := sizel + sizer
		lnk := link{size, l, r}
		ptr := vm.allocateLink(lnk)
		taggedptr := setFlag(ptr, pair)
		vm.push(taggedptr)		
	}
}

func (vm *machine) allocateLink(l link) word {
	address := vm.free
	vm.free++
	vm.links[address] = l
	return address
}

func (vm *machine) item() {
	val := vm.pop()
	if !hasFlag(val, leaf) {
		log.Fatal("item not defined for that val")
	}
	vm.push(clearFlag(val, leaf))
}

func(vm *machine) sizeOf(w word) word {
	var size word
	if hasFlag(w, pair) {
		index := clearFlag(w, pair)
		size = vm.links[index].size
	} else if hasFlag(w, leaf) {
		size = 1
	} else if hasFlag(w, nil) {
		size = 0
	} else {
		log.Fatal("size only defined on lists")
	}
	
	return size
}

