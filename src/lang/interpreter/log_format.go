package interpreter

import "fmt"
import "strings"

func log_format(in Action, hash_spacing int, endl bool) {

  switch in.Type {

    case "hash":
      if len(in.Hash_Values) == 0 {
        fmt.Print("[::]")
        goto end
      }

      fmt.Println("[:")

      for k, v := range in.Hash_Values {
        fmt.Print(strings.Repeat(" ", hash_spacing) +  k + ": ")
        log_format(v[0], hash_spacing + 2, true)
      }

      fmt.Print(strings.Repeat(" ", hash_spacing - 2) + ":]")

    case "array":
      if len(in.Hash_Values) == 0 {
        fmt.Print("[]")
        goto end
      }

      fmt.Println("[")

      for k, v := range in.Hash_Values {
        fmt.Print(strings.Repeat(" ", hash_spacing) +  k + ": ")
        log_format(v[0], hash_spacing + 2, true)
      }

      fmt.Print(strings.Repeat(" ", hash_spacing - 2) + ":]")
    case "group":
      fmt.Print("{...}")
    case "process":
      fmt.Print("{...}", "PARAM COUNT:", len(in.Params))
    case "operation":
      log_format(in.First[0], hash_spacing, false)
      fmt.Print("", getOp(in.Type), "")
      log_format(in.Second[0], hash_spacing, false)
    default:
      //cast to a string, then print

  }

  end:
  if endl { //if it was "logged", print a newline
    fmt.Println()
  }
}
