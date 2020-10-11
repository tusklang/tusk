package compiler

//list of keywords in json
var tokensJSON = []byte(
	`
[
  {
    "name": "var",
    "remove": "var",
    "pattern": "(var(\\s+))",
    "type": "id"
  },
  {
    "name": "build",
    "remove": "build",
    "pattern": "(build(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "access",
    "remove": "access",
    "pattern": "(access(\\s+))",
    "type": "id"
  },
  {
    "name": "defer",
    "remove": "defer",
    "pattern": "(defer(\\s+))",
    "type": "id"
  },
  {
    "name": "ovld",
    "remove": "ovld",
    "pattern": "(ovld(\\s+))",
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
    "pattern": "(\\-\\-)",
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
    "name": "^=",
    "remove": "^=",
    "pattern": "(\\^\\=)",
    "type": "?operation"
  },
  {
    "name": "%=",
    "remove": "%=",
    "pattern": "(\\%\\=)",
    "type": "?operation"
  },
  {
    "name": "->",
    "remove": "->",
    "pattern": "(\\-\\>)",
    "type": "operation"
  },
  {
    "name": "^",
    "remove": "**",
    "pattern": "(\\*\\*)",
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
    "pattern": "(fn(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "proto",
    "remove": "proto",
    "pattern": "(proto(\\s*)\\{)",
    "type": "id"
  },
  {
    "name": "static",
    "remove": "static",
    "pattern": "(static(\\s+))",
    "type": "id"
  },
  {
    "name": "instance",
    "remove": "instance",
    "pattern": "(instance(\\s+))",
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
    "pattern": "(return(\\s+))",
    "type": "id"
  },
  {
    "name": "if",
    "remove": "if",
    "pattern": "(if\\s*\\()",
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
    "pattern": "(elif\\s*\\()",
    "type": "id"
  },
  {
    "name": "elif",
    "remove": "elif",
    "pattern": "(elif\\s*\\{)",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "(else\\s*\\{)",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "(else\\s*\\()",
    "type": "id"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "(else\\s+)",
    "type": "id"
  },
  {
    "name": "while",
    "remove": "while",
    "pattern": "(while(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "while",
    "remove": "while",
    "pattern": "(while(\\s*)\\{)",
    "type": "id"
  },
  {
    "name": "each",
    "remove": "each",
    "pattern": "(each(\\s*)\\()",
    "type": "id"
  },
  {
    "name": "each",
    "remove": "each",
    "pattern": "(each(\\s*)\\{)",
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
    "pattern": "(break\\s*\\,)",
    "type": "loopcont"
  },
  {
    "name": "continue",
    "remove": "continue",
    "pattern": "(continue\\s*\\,)",
    "type": "loopcont"
  },
  {
    "name": "try",
    "remove": "try",
    "pattern": "(try\\s*(\\{|\\())",
    "type": "id"
  },
  {
    "name": "catch",
    "remove": "catch",
    "pattern": "(catch\\s*(\\{|\\())",
    "type": "id"
  }
]
`,
)
