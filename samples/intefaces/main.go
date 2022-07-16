package main

import (
	"fmt"
)

type CustomError struct {
}

func (c *CustomError) Error() string {
	return "custom error"
}

func GetError() error {
	var c *CustomError
	fmt.Printf("%v (%v)\n", c, c == nil)
	return c
}

func GetError1() error {
	var c error
	fmt.Printf("%v (%v)\n", c, c == nil)
	c = &CustomError{}
	return c
}

func main() {
	res := GetError()
	_ = res
	fmt.Printf("%#v: %v\n", res, res == nil)

	res = GetError1()
	fmt.Printf("%#v: %v\n", res, res == nil)
}
