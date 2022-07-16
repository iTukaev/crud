package main

import (
	"fmt"
)

func panicCall(a int) int {
	var i = (5 / 2) - a
	return 10 / i
}

func catchPanic() error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover: ", e)
		}
	}()
	//err := errors.New("hellos")
	panicCall(2)
	return nil
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover in main: ", e)
		}
	}()
	err := catchPanic()
	fmt.Println(err)
}
