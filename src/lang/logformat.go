package lang

import "os/exec"
import "os"
import "strings"

func log(val string) {
  cmd := exec.Command("./console/main.exe")

  cmd.Stdin = strings.NewReader(val + "\n")
  cmd.Stdout = os.Stdout

  _ = cmd.Run()
}
