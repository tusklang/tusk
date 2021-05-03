package tokenizer

import "github.com/dlclark/regexp2"

var tokenlist = []*regexp2.Regexp{

	/************ whitespace ************/
	regexp2.MustCompile("\\n", 0),  //newline
	regexp2.MustCompile("\\s+", 0), //whitespace
	/************************************/

	/************ keywords ************/
	regexp2.MustCompile("fn(?=(\\z|\\(|\\{|\\s+))", 0),     //fn
	regexp2.MustCompile("return(?=(\\z|\\(|\\{|\\s+))", 0), //return
	regexp2.MustCompile("var(?=\\z|\\s+)", 0),              //var
	regexp2.MustCompile("if(?=(\\z|\\(|\\{|\\s+))", 0),     //for
	regexp2.MustCompile("else(?=(\\z|\\(|\\{|\\s+))", 0),   //for
	regexp2.MustCompile("for(?=(\\z|\\(|\\{|\\s+))", 0),    //for
	regexp2.MustCompile("while(?=(\\z|\\(|\\{|\\s+))", 0),  //while
	/**********************************/

	/************ braces ************/
	regexp2.MustCompile("\\(", 0), //opening parenthesis
	regexp2.MustCompile("\\)", 0), //closing parenthesis
	regexp2.MustCompile("\\{", 0), //opening curly brace
	regexp2.MustCompile("\\}", 0), //closing curly brace
	/********************************/

	/************ operators ************/
	regexp2.MustCompile("(?<=[\\d.])\\s*\\+(?=[\\d.])", 0), // +
	regexp2.MustCompile("(?<=[\\d.])\\s*\\+(?=[\\d.])", 0), // -
	regexp2.MustCompile("\\*", 0),                          // *
	regexp2.MustCompile("\\/", 0),                          // /
	regexp2.MustCompile("\\*\\*", 0),                       // **
	/***********************************/

	/************ misc ************/
	regexp2.MustCompile("([+-]*[0-9]*\\.[0-9]*)", 0),     //floating literal
	regexp2.MustCompile("([+-]*\\d+)", 0),                //integer literal
	regexp2.MustCompile("([a-zA-Z$_][a-zA-Z$_0-9]*)", 0), //variable
	/******************************/
}
