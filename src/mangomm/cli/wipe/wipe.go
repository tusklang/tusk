package mango_wipe

import "os"

//export Wipe
func Wipe() {

  dir := os.Args[1]

  //remove the .mngo file
  os.Remove(dir + ".mngo")

  //remove the saved oats
  os.RemoveAll(dir + "/mango/")
}
