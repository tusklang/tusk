package compiler

//list of keywords in json
var keywordJSON =
[]byte(
`
[
  {
    "name": "var",
    "remove": "var",
    "pattern": "(var(\\s+))",
    "type": "id"
  },
  {
    "name": "var",
    "remove": "var",
    "pattern": "(var(\\s*)\\~)",
    "type": "id_non_tilde"
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
    "name": "log",
    "remove": "log",
    "pattern": "(log(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "log",
    "remove": "log",
    "pattern": "(log(\\s+))",
    "type": "id"
  },
  {
    "name": "print",
    "remove": "print",
    "pattern": "(print(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "print",
    "remove": "print",
    "pattern": "(print(\\s+))",
    "type": "id"
  },
  {
    "name": "->",
    "remove": "->",
    "pattern": "(\\-\\>)",
    "type": "operation"
  },
  {
    "name": "=>",
    "remove": "=>",
    "pattern": "(\\=\\>)",
    "type": "operation"
  },
  {
    "name": "^",
    "remove": "**",
    "pattern": "(\\*\\*)",
    "type": "operation"
  },
  {
    "name": "~",
    "remove": "~",
    "pattern": "\\~",
    "type": "operation"
  },
  {
    "name": "[:",
    "remove": "[:",
    "pattern": "(\\[\\:)",
    "type": "?open_brace"
  },
  {
    "name": ":]",
    "remove": ":]",
    "pattern": "(\\:\\])",
    "type": "?close_brace"
  },
  {
    "name": "::",
    "remove": "::",
    "pattern": "(\\:\\:)",
    "type": "operation"
  },
  {
    "name": ":",
    "remove": ":",
    "pattern": "\\:",
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
    "name": "<~",
    "remove": "<~",
    "pattern": "(\\<\\~)",
    "type": "operation"
  },
  {
    "name": "<-",
    "remove": "<-",
    "pattern": "(\\<\\-)",
    "type": "operation"
  },
  {
    "name": "fargc",
    "remove": "fargc",
    "pattern": "(fargc\\s+)",
    "type": "id"
  },
  {
    "name": "fargc",
    "remove": "fargc",
    "pattern": "(fargc\\s*\\~)",
    "type": "id_non_tilde"
  },
  {
    "name": "function",
    "remove": "fn",
    "pattern": "(fn(\\s*)\\()",
    "type": "id"
  },
  {
    "name": ",",
    "remove": ",",
    "pattern": "\\,",
    "type": "operation"
  },
  {
    "name": "await",
    "remove": "await",
    "pattern": "(await(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "await",
    "remove": "await",
    "pattern": "(await(\\s+))",
    "type": "id"
  },
  {
    "name": "return",
    "remove": "return",
    "pattern": "(return(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "return",
    "remove": "return",
    "pattern": "(return(\\s+))",
    "type": "id"
  },
  {
    "name": "?",
    "remove": "?",
    "pattern": "(\\?)",
    "type": "operation"
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
    "name": "!=",
    "remove": "!=",
    "pattern": "(\\!\\=)",
    "type": "operation"
  },
  {
    "name": "=",
    "remove": "==",
    "pattern": "(\\=\\=)",
    "type": "operation"
  },
  {
    "name": "=",
    "remove": "=",
    "pattern": "\\=",
    "type": "operation"
  },
  {
    "name": "include",
    "remove": "include",
    "pattern": "(include(\\s+))",
    "type": "id"
  },
  {
    "name": "include",
    "remove": "include",
    "pattern": "(include(\\s*)(~))",
    "type": "id_non_tilde"
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
    "type": "id_non_tilde"
  },
  {
    "name": "continue",
    "remove": "continue",
    "pattern": "(continue\\s*\\,)",
    "type": "id_non_tilde"
  }
]
`,
)
