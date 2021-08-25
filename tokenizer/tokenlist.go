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
	/**********************************/

	/************ braces ************/
	{("\\("), "("}, //opening parenthesis
	{("\\)"), ")"}, //closing parenthesis
	{("\\{"), "{"}, //opening curly brace
	{("\\}"), "}"}, //closing curly brace
	/********************************/

	/************ operators ************/
	{("(?<=[\\d.]\\s*)\\+(?=\\s*[\\d.])"), "+"}, // +
	{("(?<=[\\d.]\\s*)\\-(?=\\s*[\\d.])"), "-"}, // -
	{("\\*"), "*"},     // *
	{("\\/"), "/"},     // /
	{("\\*\\*"), "**"}, // **
	{("\\="), "="},     // =
	{("\\:"), ":"},     // :
	/***********************************/

	/************ types ************/
	{"int64(?![a-zA-Z$_0-9])", "dtype"}, //int64 type
	{"int32(?![a-zA-Z$_0-9])", "dtype"}, //int32 type
	{"int16(?![a-zA-Z$_0-9])", "dtype"}, //int16 type
	{"int8(?![a-zA-Z$_0-9])", "dtype"},  //int8 type
	{"int(?![a-zA-Z$_0-9])", "dtype"},   //int type

	{"uint64(?![a-zA-Z$_0-9])", "dtype"}, //uint64 type
	{"uint32(?![a-zA-Z$_0-9])", "dtype"}, //uint32 type
	{"uint16(?![a-zA-Z$_0-9])", "dtype"}, //uint16 type
	{"uint8(?![a-zA-Z$_0-9])", "dtype"},  //uint8 type
	{"uint(?![a-zA-Z$_0-9])", "dtype"},   //uint type

	{"float64(?![a-zA-Z$_0-9])", "dtype"}, //float64 type
	{"float32(?![a-zA-Z$_0-9])", "dtype"}, //float32 type
	{"float(?![a-zA-Z$_0-9])", "dtype"},   //float type

	{"bool(?![a-zA-Z$_0-9])", "dtype"},   //boolean type
	{"char(?![a-zA-Z$_0-9])", "dtype"},   //char type
	{"string(?![a-zA-Z$_0-9])", "dtype"}, //string type
	/*******************************/

	/************ misc ************/
	{"([+-]*[0-9]*\\.[0-9]*)", "float"},       //floating literal
	{"([+-]*\\d+)", "int"},                    //integer literal
	{"([a-zA-Z$_][a-zA-Z$_0-9]*)", "varname"}, //variable
	/******************************/
}
