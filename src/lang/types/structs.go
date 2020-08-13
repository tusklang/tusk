package types

//export Action
type Action struct {
  Type            string
  Name            string
  Value           OmmType
  ExpAct        []Action

  //stuff for operations

  First         []Action
  Second        []Action

  //////////////////////

  //stuff for runtime arrays and hashes

  Array       [][]Action
  Hash     [][2][]Action

  /////////////////////////////////////

  //stuff to panic errors and give stack

  File            string
  Line            uint64

  //////////////////////////////////////
}

type Returner struct {
  Exp      *OmmType
  Type      string
}

type Oat struct {
  Actions    []Action
  Variables    map[string][]Action
}
