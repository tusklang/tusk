package interpreter

//export Condition
type Condition struct {
  Type            string
  Condition     []Action
  Actions       []Action
}

type SubCaller struct {
  Indexes     [][]Action
  Args        [][]Action
  IsProc          bool
}

//export Action
type Action struct {
  Type            string
  Name            string
  ExpStr          string
  ExpAct        []Action
  Params        []string
  Args        [][]Action
  Condition     []Condition

  //stuff for operations

  First         []Action
  Second        []Action
  Degree        []Action

  //stuff for indexes

  Value       [][]Action
  Indexes     [][]Action
  Hash_Values     map[string][]Action

  IsMutable       bool
  Access          string
  SubCall       []SubCaller

  //stuff for numbers
  Integer       []int64
  Decimal       []int64
}

type Variable struct {
  Type      string
  Name      string
  Value     Action
  GoProc    func(actions []Action, cli_params CliParams, vars map[string]Variable, expReturn bool, this_vals []Action, dir string) Returner
}

type Returner struct {
  Variables map[string]Variable
  Exp       Action
  Type      string
}
