package compiler

//export Reverse
func Reverse(s string) string {
  var result = ""

  for _,v := range s {
    result = string(v) + result
  }
  return result
}

func reverseInterface(in []interface{}) {

  for i, o := 0, len(in) - 1; i < o; i, o = i + 1, o - 1 {
    in[i], in[o] = in[o], in[i]
  }
}

func reverseStringSlice(in []string) {

  for i, o := 0, len(in) - 1; i < o; i, o = i + 1, o - 1 {
    in[i], in[o] = in[o], in[i]
  }
}
