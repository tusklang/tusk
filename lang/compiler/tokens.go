package compiler

//list of keywords in json
var tokensJSON = []byte(
	`
[
  {
    "name": "var",
    "remove": "var",
    "pattern": "\\b(var(\\s+))",
    "type": "id"
  },
  {
    "name": "build",
    "remove": "build",
    "pattern": "\\b(build(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "access",
    "remove": "access",
    "pattern": "\\b(access(\\s+))",
    "type": "id"
  },
  {
    "name": "defer",
    "remove": "defer",
    "pattern": "\\b(defer(\\s+))",
    "type": "id"
  },
  {
    "name": "ovld",
    "remove": "ovld",
    "pattern": "\\b(ovld(\\s+))",
    "type": "id"
  },
  {
    "name": "!=",
    "remove": "!=",
    "pattern": "(\\!\\=)",
    "type": "operation"
  },
  {
    "name": "==",
    "remove": "==",
    "pattern": "\\=\\=",
    "type": "operation"
  },
  {
    "name": ":=",
    "remove": ":=",
    "pattern": "(\\:\\=)",
    "type": "operation"
  },
  {
    "name": "->",
    "remove": "->",
    "pattern": "(\\-\\>)",
    "type": "operation"
  },
  {
    "name": "++",
    "remove": "++",
    "pattern": "(\\+\\+)",
    "type": "?operation"
  },
  {
    "name": "+=",
    "remove": "+=",
    "pattern": "(\\+\\=)",
    "type": "?operation"
  },
  {
    "name": "--",
    "remove": "--",
    "pattern": "(\\-\\-)",
    "type": "?operation"
  },
  {
    "name": "-=",
    "remove": "-=",
    "pattern": "(\\-\\=)",
    "type": "?operation"
  },
  {
    "name": "*=",
    "remove": "*=",
    "pattern": "(\\*\\=)",
    "type": "?operation"
  },
  {
    "name": "/=",
    "remove": "/=",
    "pattern": "(\\/\\=)",
    "type": "?operation"
  },
  {
    "name": "%=",
    "remove": "%=",
    "pattern": "(\\%\\=)",
    "type": "?operation"
  },
  {
    "name": "^=",
    "remove": "^=",
    "pattern": "(\\^\\=)",
    "type": "?operation"
  },
  {
    "name": "->",
    "remove": "->",
    "pattern": "(\\-\\>)",
    "type": "operation"
  },
  {
    "name": "[:",
    "remove": "[:",
    "pattern": "\\[\\:",
    "type": "?open_brace"
  },
  {
    "name": ":]",
    "remove": ":]",
    "pattern": "\\:\\]",
    "type": "?close_brace"
  },
  {
    "name": "::",
    "remove": "::",
    "pattern": "(\\:\\:)",
    "type": "operation"
  },
  {
    "name": "=",
    "remove": "=",
    "pattern": "\\=",
    "type": "operation"
  },
  {
    "name": "+",
    "remove": "+",
    "pattern": "(\\+)",
    "type": "operation"
  },
  {
    "name": "-",
    "remove": "-",
    "pattern": "(\\-)",
    "type": "operation"
  },
  {
    "name": "*",
    "remove": "*",
    "pattern": "\\*",
    "type": "operation"
  },
  {
    "name": "/",
    "remove": "/",
    "pattern": "\\/",
    "type": "operation"
  },
  {
    "name": "%",
    "remove": "%",
    "pattern": "\\%",
    "type": "operation"
  },
  {
    "name": "^",
    "remove": "^",
    "pattern": "\\^",
    "type": "operation"
  },
  {
    "name": "(",
    "remove": "(",
    "pattern": "\\(",
    "type": "?open_brace"
  },
  {
    "name": ")",
    "remove": ")",
    "pattern": "\\)",
    "type": "?close_brace"
  },
  {
    "name": "[",
    "remove": "[",
    "pattern": "\\[",
    "type": "?open_brace"
  },
  {
    "name": "]",
    "remove": "]",
    "pattern": "\\]",
    "type": "?close_brace"
  },
  {
    "name": "{",
    "remove": "{",
    "pattern": "\\{",
    "type": "?open_brace"
  },
  {
    "name": "}",
    "remove": "}",
    "pattern": "\\}",
    "type": "?close_brace"
  },
  {
    "name": "function",
    "remove": "fn",
    "pattern": "\\b(fn(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "proto",
    "remove": "proto",
    "pattern": "\\b(proto(\\s*)\\{)",
    "type": "id"
  },
  {
    "name": "static",
    "remove": "static",
    "pattern": "\\b(static(\\s+))",
    "type": "id"
  },
  {
    "name": "instance",
    "remove": "instance",
    "pattern": "\\b(instance(\\s+))",
    "type": "id"
  },
  {
    "name": ",",
    "remove": ",",
    "pattern": "\\,",
    "type": "operation"
  },
  {
    "name": "return",
    "remove": "return",
    "pattern": "\\b(return(\\s+))",
    "type": "id"
  },
  {
    "name": "if",
    "remove": "if",
    "pattern": "\\b(if\\s*\\()",
    "type": "id"
  },
  {
    "name": "if",
    "remove": "if",
    "pattern": "(if\\s*\\{)",
    "type": "id"
  },
  {
    "name": "elif",
    "remove": "elif",
    "pattern": "\\b(elif\\s*\\()",
    "type": "id"
  },
  {
    "name": "elif",
    "remove": "elif",
    "pattern": "\\b(elif\\s*\\{)",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "\\b(else\\s*\\{)",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "\\b(else\\s*\\()",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "\\b(else\\s+)",
    "type": "id"
  },
  {
    "name": "while",
    "remove": "while",
    "pattern": "\\b(while(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "while",
    "remove": "while",
    "pattern": "\\b(while(\\s*)\\{)",
    "type": "id"
  },
  {
    "name": "each",
    "remove": "each",
    "pattern": "\\b(each(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "each",
    "remove": "each",
    "pattern": "\\b(each(\\s*)\\{)",
    "type": "id"
  },
  {
    "name": "<=",
    "remove": "<=",
    "pattern": "(\\<\\=)",
    "type": "operation"
  },
  {
    "name": ">=",
    "remove": ">=",
    "pattern": "(\\>\\=)",
    "type": "operation"
  },
  {
    "name": ">",
    "remove": ">",
    "pattern": "\\>",
    "type": "operation"
  },
  {
    "name": "<",
    "remove": "<",
    "pattern": "\\<",
    "type": "operation"
  },
  {
    "name": ":",
    "remove": ":",
    "pattern": "\\:",
    "type": "operation"
  },
  {
    "name": "?",
    "remove": "?",
    "pattern": "\\?",
    "type": "operation"
  },
  {
    "name": "!",
    "remove": "!",
    "pattern": "\\!",
    "type": "?operation"
  },
  {
    "name": "&",
    "remove": "&&",
    "pattern": "(\\&\\&)",
    "type": "operation"
  },
  {
    "name": "|",
    "remove": "||",
    "pattern": "(\\|\\|)",
    "type": "operation"
  },
  {
    "name": "&",
    "remove": "&",
    "pattern": "\\&",
    "type": "operation"
  },
  {
    "name": "|",
    "remove": "|",
    "pattern": "\\|",
    "type": "operation"
  },
  {
    "name": "break",
    "remove": "break",
    "pattern": "\\b(break)\\b",
    "type": "loopcont"
  },
  {
    "name": "continue",
    "remove": "continue",
    "pattern": "\\b(continue)\\b",
    "type": "loopcont"
  },
  {
    "name": "try",
    "remove": "try",
    "pattern": "\\b(try\\s*(\\{|\\())",
    "type": "id"
  },
  {
    "name": "catch",
    "remove": "catch",
    "pattern": "\\b(catch\\s*(\\{|\\())",
    "type": "id"
  },
  {
    "name": "true",
    "remove": "true",
    "pattern": "\\btrue\\b",
    "type": "expression value"
  },
  {
    "name": "false",
    "remove": "false",
    "pattern": "\\bfalse\\b",
    "type": "expression value"
  },
  {
    "name": "undef",
    "remove": "undef",
    "pattern": "\\bundef\\b",
    "type": "expression value"
  }
]
`,
)
