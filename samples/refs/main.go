package main

import (
	"fmt"
)

func main() {
	var i = 10
	var r *int
	var n *int64
	println(i)
	println(r)
	r = &i
	k := int64(i)
	n = &k
	fmt.Printf("-- %#v, %T\n", *n, *n)
}
