package main

func arrayContain(arr []string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
      return true
    }
  }

  return false;
}

func arrayContain2Nest(arr [][]string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arrayContain(arr[i], sub) {
      return true
    }
  }

  return false
}

func indexOf2Nest(sub string, arr [][]string) []int {
  for i := 0; i < len(arr); i++ {
    for o := 0; o < len(arr[i]); o++ {
      if arr[i][o] == sub {
        return []int{ i, o }
      }
    }
  }

  return []int{ -1, -1 }
}

func RepeatAdd(s string, times int) string {
  returner := ""

  for ;times > 0; times-- {
    returner+=s
  }

  return returner
}
