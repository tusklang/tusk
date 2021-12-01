package tokenizer

type TokenItem struct {
	regexp    string
	tokentype string
}

var keywords = []string{"fn", "return", "var", "if", "else", "while", "for", "pub", "prv", "ptr", "stat", "link", "construct", "this"}

var tokenlist = []TokenItem{

	{("\\;"), "terminator"}, //semicolon
	{("\\,"), "terminator"}, //comma

	/************ comments ************/
	{"\\/\\/.*?(?=\\n|$)", "comment"},     //single line comment
	{"\\/\\*[\\s\\S]*?\\*\\/", "comment"}, //multi line comment https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch07s06.html#:~:text=We%20use%20%E2%80%B9%20.,last%20*/%20in%20the%20file.
	/**********************************/

	/************ whitespace ************/
	{("\\n"), "newline"},       //newline
	{("\\s{1}"), "whitespace"}, //whitespace
	/************************************/

	/************ braces ************/
	{("\\("), "("}, //opening parenthesis
	{("\\)"), ")"}, //closing parenthesis
	{("\\{"), "{"}, //opening curly brace
	{("\\}"), "}"}, //closing curly brace
	{("\\["), "["}, //opening square brace
	{("\\]"), "]"}, //closing square brace
	/********************************/

	/************ operators ************/
	{("\\-\\>"), "operation"}, // ->
	{("\\@"), "operation"},    // @
	{("\\#"), "operation"},    // #
	{("\\+"), "operation"},    // +
	{("\\-"), "operation"},    // -
	{("\\*"), "operation"},    // *
	{("\\/"), "operation"},    // /
	{("\\!\\="), "operation"}, // !=
	{("\\=\\="), "operation"}, // ==
	{("\\>\\="), "operation"}, // >=
	{("\\>"), "operation"},    // >
	{("\\<\\="), "operation"}, // <=
	{("\\<"), "operation"},    // <
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
