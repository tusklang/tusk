package tokenizer

type TokenItem struct {
	regexp    string
	tokentype string
}

var keywords = []string{"fn", "return", "var", "if", "else", "while", "pub", "prv", "ptr", "stat", "link", "construct", "this"}

var tokenlist = []TokenItem{

	{("\\;"), "terminator"}, //semicolon
	{("\\,"), "terminator"}, //comma

	/************ whitespace ************/
	{("\\n"), "newline"},       //newline
	{("\\s{1}"), "whitespace"}, //whitespace
	/************************************/

	/************ braces ************/
	{("\\("), "("}, //opening parenthesis
	{("\\)"), ")"}, //closing parenthesis
	{("\\{"), "{"}, //opening curly brace
	{("\\}"), "}"}, //closing curly brace
	/********************************/

	/************ operators ************/
	{("\\-\\>"), "operation"}, // ->
	{("\\+"), "operation"},    // +
	{("\\-"), "operation"},    // -
	{("\\*"), "operation"},    // *
	{("\\/"), "operation"},    // /
	{("\\=\\="), "operation"}, // ==
	{("\\>"), "operation"},    // >
	{("\\>\\="), "operation"}, // >=
	{("\\<"), "operation"},    // <
	{("\\<\\="), "operation"}, // <=
	{("\\="), "operation"},    // =
	{("\\:"), "operation"},    // :
	{("\\."), "operation"},    // .
	{("\\~"), "operation"},    // ~
	{("\\&"), "operation"},    // &
	{("\\|"), "operation"},    // |
	{("\\^"), "operation"},    // ^
	/***********************************/

	/************ misc ************/
	{"null(?![a-zA-Z$_0-9])", "null"},                          //null value
	{"([\"])((\\\\{2})*|(.*?[^\\\\](\\\\{2})*))\\1", "string"}, //string value https://stackoverflow.com/a/17231632/10696946
	{"([+-]*[0-9]*\\.[0-9]*)", "float"},                        //floating literal
	{"([+-]*\\d+)", "int"},                                     //integer literal
	{"([a-zA-Z$_][a-zA-Z$_0-9]*)", "varname"},                  //variable
	/******************************/
}
