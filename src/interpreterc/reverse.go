package main

func Reverse(s string) string {
  var result = ""
  
  for _,v := range s {
    result = string(v) + result
  }
  return result
}
