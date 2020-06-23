package interpreter

//list of commonly used values
var undef = Action{ "falsey", "exp_value", "undef", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, make(chan Returner) }
var arr = Action{ "array", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, make(chan Returner) }
var hash = Action{ "hash", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, make(chan Returner) }
var zero = Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{ 0 }, []int64{}, make(chan Returner) }
var thread = Action{ "thread", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, make(chan Returner) }
//////////////////////////////
