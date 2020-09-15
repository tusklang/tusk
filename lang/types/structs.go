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
	Value  KaType
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
	Exp  *KaType
	Type string
}

type Kast struct {
	Actions   []Action
	Variables map[string][]Action
}
