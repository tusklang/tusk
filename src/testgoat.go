package main

import "lang/interpreter"

func main() {
  var instance interpreter.Instance
  instance.FromOat("test.oat")
  instance.CallFunc("main")
}
