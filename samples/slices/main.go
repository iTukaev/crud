package main

import (
	"fmt"
)

func printSlice(a []int, name string) {
	fmt.Printf("<%v>\tlen: <%2v>, cap: <%2v>, val: <%#v>\n", name, len(a), cap(a), a)
}

func main() {
	var a []int
	if a == nil {
		fmt.Println("is nil")
	}

	for i := 0; i < 10; i++ {
		a = append(a, i)
		//printSlice(a, "a")
	}
	printSlice(a, "a")

	b := make([]int, 3)
	copy(b, a[2:5])
	printSlice(b, "b")
	b[1] = 100
	b = append(b, -1)
	printSlice(a, "a")
	printSlice(b, "b")

	c := make([]int, len(a[6:]))
	printSlice(c, "c")
	c = append(c, -2)
	printSlice(c, "c")
	printSlice(a, "a")
	b = append([]int{}, a[6:]...)
	printSlice(b, "b")

}
