package types

//export Action
type Action struct {
  Type            string              `oat:"t"`
  Name            string              `oat:"n"`
  Value           OmmType             `oat:"v"`
  ExpAct        []Action              `oat:"e"`

  //stuff for operations

  First         []Action              `oat:"f"`
  Second        []Action              `oat:"s"`

  //////////////////////

  //stuff for runtime arrays and hashes

  Array       [][]Action              `oat:"a"`
  Hash            map[string][]Action `oat:"h"`

  /////////////////////////////////////

  //stuff to panic errors and give stack

  File            string              `oat:"f"`
  Line            uint64              `oat:"l"`

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
