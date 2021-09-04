package tokenizer

type TokenItem struct {
	regexp    string
	tokentype string
}

var tokenlist = []TokenItem{

	{("\\;"), "terminator"}, //semicolon

	/************ whitespace ************/
	{("\\n"), "newline"},     //newline
	{("\\s+"), "whitespace"}, //whitespace
	/************************************/

	/************ keywords ************/
	{("fn(?=(\\z|\\(|\\{|\\s+))"), "fn"},         //fn
	{("return(?=(\\z|\\(|\\{|\\s+))"), "return"}, //return
	{("var(?=\\z|\\s+)"), "var"},                 //var
	{("if(?=(\\z|\\(|\\{|\\s+))"), "if"},         //if
	{("else(?=(\\z|\\(|\\{|\\s+))"), "else"},     //else
	{("while(?=(\\z|\\(|\\{|\\s+))"), "while"},   //while
	{("pub(?![a-zA-Z$_0-9])"), "pub"},            //pub
	{("prt(?![a-zA-Z$_0-9])"), "prt"},            //prt
	{("stat(?![a-zA-Z$_0-9])"), "stat"},          //stat
	/**********************************/

	/************ braces ************/
	{("\\("), "("}, //opening parenthesis
	{("\\)"), ")"}, //closing parenthesis
	{("\\{"), "{"}, //opening curly brace
	{("\\}"), "}"}, //closing curly brace
	/********************************/

	/************ operators ************/
	{("(?<=[\\d.]\\s*)\\+(?=\\s*[\\d.])"), "operation"}, // +
	{("(?<=[\\d.]\\s*)\\-(?=\\s*[\\d.])"), "operation"}, // -
	{("\\*"), "operation"},                              // *
	{("\\/"), "operation"},                              // /
	{("\\*\\*"), "operation"},                           // **
	{("\\="), "operation"},                              // =
	{("\\:"), "operation"},                              // :
	{("\\."), "operation"},                              // .
	/***********************************/

	/************ types ************/
	{"i64(?![a-zA-Z$_0-9])", "dtype"}, //int64 type
	{"i32(?![a-zA-Z$_0-9])", "dtype"}, //int32 type
	{"i16(?![a-zA-Z$_0-9])", "dtype"}, //int16 type
	{"i8(?![a-zA-Z$_0-9])", "dtype"},  //int8 type

	{"u64(?![a-zA-Z$_0-9])", "dtype"}, //uint64 type
	{"u32(?![a-zA-Z$_0-9])", "dtype"}, //uint32 type
	{"u16(?![a-zA-Z$_0-9])", "dtype"}, //uint16 type
	{"u8(?![a-zA-Z$_0-9])", "dtype"},  //uint8 type

	{"f64(?![a-zA-Z$_0-9])", "dtype"}, //float64 type
	{"f32(?![a-zA-Z$_0-9])", "dtype"}, //float32 type

	{"bool(?![a-zA-Z$_0-9])", "dtype"},   //boolean type
	{"string(?![a-zA-Z$_0-9])", "dtype"}, //string type
	/*******************************/

	/************ misc ************/
	{"([+-]*[0-9]*\\.[0-9]*)", "float"},       //floating literal
	{"([+-]*\\d+)", "int"},                    //integer literal
	{"([a-zA-Z$_][a-zA-Z$_0-9]*)", "varname"}, //variable
	/******************************/
}
