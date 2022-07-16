package main

import "fmt"

func callIt(a int) {
	defer println("first defer:", a)
	defer func() {
		fmt.Printf("second defer: %v\n", a)
	}()

	a += 10
}

func main() {
	callIt(100)
}
