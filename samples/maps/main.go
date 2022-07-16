package main

import (
	"fmt"
	"strconv"
)

func printMap(a map[string]int) {
	fmt.Printf("len: <%2d>, <%+v>\n", len(a), a)
}

func fixIt(a map[string]int) {
	a["1"] = 100
	a["ups"] = -1
}

func main() {
	a := map[string]int{}
	for i := 0; i < 3; i++ {
		a[strconv.FormatInt(int64(i), 10)] = i
	}
	printMap(a)
	fixIt(a)
	printMap(a)
}
