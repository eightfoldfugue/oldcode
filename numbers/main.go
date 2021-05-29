package main

import (
	"fmt"
)


func main() {
	x := []byte{255}
 	y := []byte{3}

 	z := AddBigs(x, y)


 	fmt.Println(x)
  	fmt.Println(y)
 	fmt.Println(z)

}


// notice that this is reverse endian of normal, big or little idk
// in a slimmer implementation, output array is at most
// one more than the largest input array.
// if the final carry is one we fill the last spot and bump
// the length,
// otherwise the final spot is empty and we don't bump
// if Go has some kind of !SETLENGTH option on slices we could
// do this now...
func AddBigs(x, y []byte) (z []byte) {
	var c byte = 0
	for i := 0; i < Max(len(x), len(y)); i++ {
		z_i, c_next := AddBytes(x[i], y[i], c)
		z = append(z, z_i)
		c = c_next
	}
	if c == 1 {
		z = append(z, c)
	}
	return
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func AddBytes(x, y, cIn byte) (sum, cOut byte) {
	sum = x + y + cIn
	cOut = ((x & y) | ((x | y) &^ sum)) >> 7
	return
}

func Show(x byte) {
	fmt.Printf("%.8b\n", x)
}


/*

z := []
c := 0
for i := 0, i < Max(len(x), len(y)); i ++ {
	z_i, c_next := Add(x, y, c)
	z = append(z, z_i)
	c = c_next
}
if c == 1 {
	z = append(z, c)
}



if we define append in such a way that 
append bignum 0 = bignum

then it is always the case that we append the carry at then end...

but! now we check on every iteration to see if the thing to be
added is 0. elegent in form, wasteful in process



*/
