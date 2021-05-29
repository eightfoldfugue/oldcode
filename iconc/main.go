package main

import (
	"fmt"
)

type List interface{
	Size() int
	Level() int
}

type empty struct{}
func (e empty) Size() int {return 0}
func (e empty) Level() int {return 0}

type single struct {
	item interface{}
}
func (s single) Size() int {return 1}
func (s single) Level() int {return 0}

type branch struct {
	size  int
	level int
	left List
	right List
}
func (b branch) Size() int {return b.size}
func (b branch) Level() int {return b.level}



func isEmpty(xs List) bool {
	_, ok := xs.(empty)
	return ok
}
func isSingle(xs List) bool {
	_, ok := xs.(single)
	return ok
}
func Split(xs List) (List, List) {
	ys := xs.(*branch).left
	zs := xs.(*branch).right
	return ys, zs
}

// why doesn't this respect the no emptys invariant?
func link(l, r List) List {
	size := l.Size() + r.Size()
	level := 1 + max(l.Level(), r.Level())
	return &branch{size, level, l, r}
}

func max(x, y int) int {
	if x > y {return x}
	return y
}

func Index(xs List, i int) interface{} {
	if i < 0 || i >= xs.Size() {
		panic("index out of range")
	}
	return indexRec(xs, i)
}

func indexRec(xs List, i int) interface{} {
	if isSingle(xs) {
		return xs.(single).item
	}
	ys, zs := Split(xs)
	if i < ys.Size() {
		return indexRec(ys, i)
	} else {
		return indexRec(zs, i - ys.Size())
	}
}

func Show(xs List) {
	if isEmpty(xs) {
		fmt.Println("[]")
	}
	fmt.Print("[")
	last := xs.Size() - 1
	for i := 0; i < last; i++ {
		fmt.Printf("%v, ", Index(xs, i))
	}
	fmt.Printf("%v]\n", Index(xs, last))
}

func update(xs List, i int, y interface{}) List {
	if isSingle(xs) {
		return single{y}
	}
	ys, zs := Split(xs)
	if i < ys.Size() {
		return link(update(ys, i, y), zs)
	} else {
		new_i := i - ys.Size()
		return link(ys, update(zs, new_i, y))
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Conc(xs, ys List) List {
	diff := ys.Level() - xs.Level()
	if abs(diff) <= 1 {
		return link(xs, ys)
	}
	return nil
}

func main() {
	a := single{"a"}
	b := single{"b"}
	c := single{"c"}

	
	list := link(a, link(b, c))

	Show(list)



	fmt.Println(abs(1))
	fmt.Println(abs(0))
	fmt.Println(abs(-7))
}
