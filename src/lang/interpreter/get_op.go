package interpreter

func getOp(val string) string {
  switch (val) {
    case "add":
      return C.CString("+")
    case "subtract":
      return C.CString("-")
    case "multiply":
      return C.CString("*")
    case "divide":
      return C.CString("/")
    case "exponentiate":
      return C.CString("^")
    case "modulo":
      return C.CString("%")
  }

  //if no operation was detected, return ?OP?
  return "?OP?"
}
