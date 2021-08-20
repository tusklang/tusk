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
	{("for(?=(\\z|\\(|\\{|\\s+))"), "for"},       //for
	{("while(?=(\\z|\\(|\\{|\\s+))"), "while"},   //while
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
	{"int64", "dtype"}, //int64 type
	{"int32", "dtype"}, //int32 type
	{"int16", "dtype"}, //int16 type
	{"int8", "dtype"},  //int8 type
	{"int", "dtype"},   //int type

	{"uint64", "dtype"}, //uint64 type
	{"uint32", "dtype"}, //uint32 type
	{"uint16", "dtype"}, //uint16 type
	{"uint8", "dtype"},  //uint8 type
	{"uint", "dtype"},   //uint type

	{"float64", "dtype"}, //float64 type
	{"float32", "dtype"}, //float32 type
	{"float", "dtype"},   //float type

	{"bool", "dtype"},   //boolean type
	{"char", "dtype"},   //char type
	{"string", "dtype"}, //string type
	/*******************************/

	/************ misc ************/
	{"([+-]*[0-9]*\\.[0-9]*)", "float"},       //floating literal
	{"([+-]*\\d+)", "int"},                    //integer literal
	{"([a-zA-Z$_][a-zA-Z$_0-9]*)", "varname"}, //variable
	/******************************/
}
