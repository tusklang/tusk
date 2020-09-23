package main

import (
	"fmt"
	"plugin"
)

func main() {
	t, e := plugin.Open("test.so")
	fmt.Println(e)
	p, e := t.Lookup("Test")
	fmt.Println(e)
	p.(func(int))(1)
}
