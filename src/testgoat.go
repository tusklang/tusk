package main

import "goat"

func main() {
  test1 := goat.CompileFile("test.omm")
  test2 := goat.NewInstance(test1)
  test2.CallFunc("main")
}
