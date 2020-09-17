package types

//export Action
type Action struct {

	//stuff to panic errors and give stack

	File string
	Line uint64

	//////////////////////////////////////

	//stuff for regular actions

	Type   string
	Name   string
	Value  TuskType
	ExpAct []Action

	///////////////////////////

	//stuff for operations

	First  []Action
	Second []Action

	//////////////////////

	//stuff for runtime arrays and hashes

	Array [][]Action
	Hash  [][2][]Action

	/////////////////////////////////////

}

type Returner struct {
	Exp  *TuskType
	Type string
}

type Tuskst struct {
	Actions   []Action
	Variables map[string][]Action
}
