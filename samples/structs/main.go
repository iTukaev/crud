package main

import (
	"fmt"
)

type MyInt int

// Написать пример
type HandlerFunc func([]byte) (int, error)

type MyCoolStruct struct {
	privateField int
	PublicField  string
}

func (m MyInt) Add(i int) MyInt {
	m = m + MyInt(i)
	return m
}

func (m *MyInt) Summ(i int) MyInt {
	*m += MyInt(i)
	return *m
}

func main() {
	mi := MyInt(5)
	fmt.Printf("Add mi before: <%v>, res: <%v>, mi after: <%v>\n", mi, mi.Add(3), mi)
	fmt.Printf("Sum mi before: <%v>, res: <%v>, mi after: <%v>\n", mi, mi.Summ(10), mi)
	m := MyCoolStruct{
		privateField: 1,
		PublicField:  "hello",
	}

	var id IncData
	v, err := id.check(&m)
	if err != nil {
		fmt.Printf("%#v\n", err)
	} else {
		fmt.Println(v)
	}
}
