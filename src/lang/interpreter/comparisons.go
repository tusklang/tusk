package interpreter

func isTruthy(val Action) bool {
  return !(val.ExpStr == "false" || val.Type == "falsey")
}
