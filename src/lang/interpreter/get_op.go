package interpreter

func getOp(val string) string {
  switch (val) {
    case "add":
      return "+"
    case "subtract":
      return "-"
    case "multiply":
      return "*"
    case "divide":
      return "/"
    case "exponentiate":
      return "^"
    case "modulo":
      return "%"
  }

  //if no operation was detected, return ?OP?
  return "?OP?"
}
