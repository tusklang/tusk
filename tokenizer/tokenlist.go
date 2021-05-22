package tokenizer

var (
	FloatPat = "([+-]*[0-9]*\\.[0-9]*)"
	IntPat   = "([+-]*\\d+)"
	VarPat   = "([a-zA-Z$_][a-zA-Z$_0-9]*)"
)

var tokenlist = []string{

	("\\;"), //semicolon

	/************ whitespace ************/
	("\\n"),  //newline
	("\\s+"), //whitespace
	/************************************/

	/************ keywords ************/
	("fn(?=(\\z|\\(|\\{|\\s+))"),     //fn
	("return(?=(\\z|\\(|\\{|\\s+))"), //return
	("var(?=\\z|\\s+)"),              //var
	("if(?=(\\z|\\(|\\{|\\s+))"),     //for
	("else(?=(\\z|\\(|\\{|\\s+))"),   //for
	("for(?=(\\z|\\(|\\{|\\s+))"),    //for
	("while(?=(\\z|\\(|\\{|\\s+))"),  //while
	/**********************************/

	/************ braces ************/
	("\\("), //opening parenthesis
	("\\)"), //closing parenthesis
	("\\{"), //opening curly brace
	("\\}"), //closing curly brace
	/********************************/

	/************ operators ************/
	("(?<=[\\d.]\\s*)\\+(?=\\s*[\\d.])"), // +
	("(?<=[\\d.]\\s*)\\-(?=\\s*[\\d.])"), // -
	("\\*"),                              // *
	("\\/"),                              // /
	("\\*\\*"),                           // **
	/***********************************/

	/************ misc ************/
	FloatPat, //floating literal
	IntPat,   //integer literal
	VarPat,   //variable
	/******************************/
}
