package compiler

//list of keywords in json
var keywordJSON =
[]byte(
`
[
  {
    "name": "local",
    "remove": "local",
    "pattern": "(local(\\s*)(~?))",
    "type": "id"
  },
  {
    "name": "local",
    "remove": "lcl",
    "pattern": "(lcl(\\s*)(~?))",
    "type": "id"
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
    "name": "~~~",
    "remove": "~~~",
    "pattern": "((~~~))",
    "type": "operation"
  },
  {
    "name": "~~",
    "remove": "~~",
    "pattern": "(~~)",
    "type": "operation"
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
    "name": "^",
    "remove": "**",
    "pattern": "(\\*\\*)",
    "type": "operation"
  },
  {
    "name": "newlineN",
    "remove": "\\n",
    "pattern": "\\n",
    "type": "newline"
  },
  {
    "name": "newlineN",
    "remove": "\\r\\n",
    "pattern": "(\\r\\n)",
    "type": "newline"
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
    "name": "async",
    "remove": "async",
    "pattern": "(async\\s*\\()",
    "type": "operation"
  },
  {
    "name": "sync",
    "remove": "sync",
    "pattern": "(sync\\s*\\()",
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
    "type": "id_non_tilde"
  },
  {
    "name": "global",
    "remove": "gbl",
    "pattern": "(gbl(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "global",
    "remove": "gbl",
    "pattern": "(gbl(\\s+))",
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
    "name": "cond",
    "remove": "cond",
    "pattern": "(cond\\s*\\[)",
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
    "name": "import",
    "remove": "import",
    "pattern": "(import(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "import",
    "remove": "import",
    "pattern": "(import(\\s+))",
    "type": "id"
  },
  {
    "name": "include",
    "remove": "include",
    "pattern": "(include(\\s+))",
    "type": "id_non_tilde"
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
    "pattern": "(break(\\s+))",
    "type": "id_non_tilde"
  },
  {
    "name": "skip",
    "remove": "skip",
    "pattern": "(skip(\\s+))",
    "type": "id_non_tilde"
  },
  {
    "name": "number",
    "remove": "number",
    "pattern": "(number\\-\\>)",
    "type": "type"
  },
  {
    "name": "number",
    "remove": "num",
    "pattern": "(num\\-\\>)",
    "type": "type"
  },
  {
    "name": "string",
    "remove": "string",
    "pattern": "(string\\-\\>)",
    "type": "type"
  },
  {
    "name": "boolean",
    "remove": "bool",
    "pattern": "(bool\\-\\>)",
    "type": "type"
  },
  {
    "name": "falsey",
    "remove": "falsey",
    "pattern": "(falsey\\-\\>)",
    "type": "type"
  },
  {
    "name": "hash",
    "remove": "hash",
    "pattern": "(hash\\-\\>)",
    "type": "type"
  },
  {
    "name": "array",
    "remove": "array",
    "pattern": "(array\\-\\>)",
    "type": "type"
  }
]
`,
)
