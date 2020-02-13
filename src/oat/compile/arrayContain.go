package main

func arrayContain(arr []string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
      return true
    }
  }

  return false;
}
