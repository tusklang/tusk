package lang

//function to chunk a string into pieces
func Chunk(val string, by int) []string {

  //val to return
  nVal := []string{}
  last := 0

  for i := 0; i < len(val); i++ {

    //if it is on the loop then append the last to the current index to nVar
    if i % by == 0 {
      nVal = append(nVal, val[last:i])
      last = i
    }
  }

  //if val's length does not equal the last value append the remaining bits of the string to nVal
  if (last != len(val)) {
    nVal = append(nVal, val[last:len(val)])
  }

  return nVal;
}
