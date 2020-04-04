package main

//export GetActType
func GetActNum(val string) int {

  switch val {
    case "string":
      return 38
    case "boolean":
      return 40
    case "variable":
      return 43
    case "array":
      return 24
    case "hash":
      return 22
    case "statement":
      return 0
    case "falsey":
      return 41
    case "none":
      return 42
    case "type":
      return 44
  }

  return 42
}
