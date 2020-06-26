package interpreter

//list of commonly used values
var undef = Action{ "falsey", "exp_value", "undef", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var arr = Action{ "array", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var hash = Action{ "hash", "hash_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var zero = Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{ 0 }, []int64{}, OmmThread{} }
var one = Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{ 1 }, []int64{}, OmmThread{} }
var neg_one = Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{ -1 }, []int64{}, OmmThread{} }
var falseAct = Action{ "boolean", "exp_value", "false", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var trueAct = Action{ "boolean", "exp_value", "true", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var emptyString = Action{ "string", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var emptyRune = Action{ "rune", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
var thread = Action{ "thread", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
//////////////////////////////
