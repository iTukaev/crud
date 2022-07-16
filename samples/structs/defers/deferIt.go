package main

import "fmt"

type TestDef struct {
	A int
}

func testDefer(a TestDef) {
	defer println("first test: ", a.A)
	defer fmt.Printf("second test: %v\n", a)
	a.A += 10
}

func testPrtDefer(a *TestDef) {
	defer println("first testPtr: ", a.A)
	defer fmt.Printf("second testPtr: %v\n", a)
	a.A += 10
}

func main() {
	a := TestDef{
		A: 8,
	}
	testDefer(a)
	//testPrtDefer(&a)
	//fmt.Printf("%#v\n", a)
}
