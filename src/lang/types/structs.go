package types

//export Action
type Action struct {
  Type            string
  Name            string
  Value           OmmType
  ExpAct        []Action

  Array       [][]Action
  Hash            map[string][]Action

  //stuff for operations

  First         []Action
  Second        []Action
  Degree        []Action

  //////////////////////

  //for compiling structures
  Static          map[string][]Action
  Instance        map[string][]Action
  //////////////////////////

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
