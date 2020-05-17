fs = require('fs');

var stdinBuffer = fs.readFileSync(0)
, { f, dir, name } = JSON.parse(stdinBuffer.toString());

global.DIRNAME = dir;
global.NAME = name;
global.warnings = [],
global.errors = [];

const
  testkey = require('./testkey'), //file to test each keyword
  processes = require('./processes'), //file with process utils
  include_parser = require('./includes'), //file that will include omm files within other omm files
  indexes = require('./indexes'), //file to allow user to write :: instead of .['name']
  id_init = require('./id_init'); //file to add ~ after every id

global.KEYWORDS = require('./keywords.json');
global.MAX_CUR_EXP = 12;
global.lexer = (file, dir) => {

  //current expression
  var curExp = ''
  //current lexxed value
  , lex = []
  //current line
  , line = 1;

  //loop through file
  outer:
  for (let i = 0; i < file.length; i++) {

    //detect a comment
    //single line comments are written as //
    if (file.substr(i).trim().startsWith('//')) {

      var end = file.substr(i).indexOf('\n');

      //if there are no more newlines, the lex is complete
      if (end == -1) break;

      i+=end;
      line++;
      continue;
    }

    while (/ |\t/.test(file[i])) i++;

    while (curExp.length > MAX_CUR_EXP) curExp = curExp.substr(1);
    while (curExp.includes('\n')) curExp = curExp.substr(curExp.indexOf('\n') + 1);

    for (let o = 0; o < KEYWORDS.length; o++) {

      //if the current key is detected
      if (testkey(KEYWORDS[o], file, i)) {

        //if it is a newline, increment "line"
        if (KEYWORDS[o].name == 'newlineN') line++;

        curExp+=KEYWORDS[o].remove;

        lex.push({
          Name: KEYWORDS[o].name,
          Exp: curExp,
          Line: line,
          Type: KEYWORDS[o].type,
          OName: KEYWORDS[o].remove,
          Dir: dir
        });

        if (KEYWORDS[o].name != 'newlineN') i+=KEYWORDS[o].remove.length - 1;

        continue outer;
      }
    }

    var substrfile = file.substr(i).trim();

    //detect a string
    if (substrfile.startsWith('\"') || substrfile.startsWith('\'') || substrfile.startsWith('\`')) {

      //get quote type
      let qType = substrfile[0],
        value = '',
        escaped = false;

      for (let o = 1; o < substrfile.length; o++, i++) {

        //if it is escape character, set escaped to true
        if (!escaped && substrfile[o] == '\\') {
          escaped = true;
          continue;
        }

        if (!escaped && substrfile[o] == qType) break;

        value+=substrfile[o];
        escaped = false;
      }

      curExp+=value;
      line+=value.match(/\n/g) == null ? 0 : value.match(/\n/g).length;

      i++;
      lex.push({
        Name: '\'' + value + '\'',
        Exp: curExp,
        Line: line,
        Type: 'string',
        OName: '\'' + value + '\'',
        Dir: dir
      });
      continue outer;
    } else if (/^(\d|\+|\-|\.)/.test(substrfile)) {

      var sign = true;

      //detect positive and negative
      while (file.substr(i).trim()[0] == '+' || file.substr(i).trim()[0] == '-')
        if (file.substr(i).trim() == '+') i++;
        else {
          sign = !sign;
          i++;
        }

      var num = '';

      if (!(/^(\d|\.)/.test(file.substr(i).trim()))) {
        if (sign) num = '1';
        else num = '-1';

        curExp+=num;

        i++;

        lex.push({
          Name: num,
          Exp: curExp,
          Line: line,
          Type: 'number',
          OName: num,
          Dir: dir
        });
        continue outer;
      }

      var substrf = file.substr(i).trim();

      while (/^(\d|\.)/.test(substrf)) {

        if (file[i] == ' ') {
          i++;
          continue;
        }

        num+=file[i];
        i++;

        substrf = file.substr(i);

        if (i > file.length) break;
      }

      curExp+=num;

      i--;
      lex.push({
        Name: (sign ? '' : '-') + num,
        Exp: curExp,
        Line: line,
        Type: 'number',
        OName: num,
        Dir: dir
      });
    } else {

      if (/\s/.test(file[i])) continue;

      var variable = '';

      var_loop:
      for (let o = i; o < file.length; o++) {

        for (let j = 0; j < KEYWORDS.length; j++)
          if (testkey(KEYWORDS[j], file, o) || /\s/.test(file[o])) break var_loop;

        variable+=file[o];
        i++;
      }

      curExp+=variable;

      variable = variable.trim();

      i--;
      lex.push({
        Name: '$' + variable,
        Exp: curExp,
        Line: line,
        Type: 'variable',
        OName: variable,
        Dir: dir
      });
    }

  }

  //account for: true, false, undefined, and null
  lex = lex.map(l => {
    if (/\$true|\$false|\$undef|\$null/.test(l.Name)) return {
      ...l,
      Name: l.Name.substr(1)
    };
    else return l;
  });

  //determine if there is an error
  //an error will occur when two adjacent tokens have the same type
  lex.forEach((v, i) => {

    if (lex[i - 1] && v.Type[0] != '?' && v.Type != 'newline') {
      if (v.Type == lex[i - 1].Type) {
        global.errors.push({
          Error:`\nUnexpected token: \"${v.OName.trim()}\" on line ${v.Line}`
          +
          `\n\nFound near: ${v.Exp.trim()}`
          +
          `\n            ${' '.repeat(v.Exp.indexOf(v.OName) == -1 ? '' : v.Exp.indexOf(v.OName))}${'^'.repeat(v.OName.length)}`
          +
          `\nAdvanced Info: Two tokens of the type \"${v.Type}\" were detected adjacent to each other`,
          Dir: dir
        });
      }
    }

  });

  lex = lex.filter(l => l.Type != "newline");

  return id_init(indexes(include_parser(global.DIRNAME, lex)));
}

console.log(
  JSON.stringify({
    WARNS: warnings,
    ERRORS: errors,
    LEX: processes.hash_inserter(lexer(processes.init(f), dir + name))
  }, null, 2)
);
