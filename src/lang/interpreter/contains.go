package interpreter

import "strings"

func containsParams(slice []string, sub string) bool {
  for _, v := range slice {
    if strings.HasPrefix(v, sub) {
      return true
    }
  }
  return false
}
