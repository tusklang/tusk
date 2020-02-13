package main

func indexOf(sub string, data []string) int {
  for k, v := range data {
    if sub == v {
      return k
    }
  }
  return -1
}
