package interpreter

//list of commonly used values
var undef = Action{ "falsey", "exp_value", "undef", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{} }
var arr = Action{ "array", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{} }
var hash = Action{ "hash", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{} }
var zero = Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{ -9, 9 }, []int64{} }
//////////////////////////////
