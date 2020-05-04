const
  testkey = require('./testkey'), //file to test each keyword
  processes = require('./processes'), //file with process utils
  include_parser = require('./includes'), //file that will include omm files within other omm files
  rw = require('./remove_whitespace') //file to remove whitespace from a string
  fs = require('fs');

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
    //single line comments are written as !>
    if (file.substr(i).startsWith('!>')) {

      var end = file.substr(i).indexOf('\n');

      //if there are no more newlines, the lex is complete
      if (end == -1) break;

      i+=end;
      line++;
      continue;
    }

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

    var substrfile = file.substr(i);

    //detect a string
    if (substrfile.startsWith('\"') || substrfile.startsWith('\'') || substrfile.startsWith('\`')) {

      //get quote type
      let qType = substrfile[0],
        value = '',
        escaped = false;

      for (let o = i + 1; o < file.length; o++, i++) {

        //if it is escape character, set escaped to true
        if (!escaped && file[o] == '\\') {
          escaped = true;
          o--;
          continue;
        }

        if (!escaped && file.substr(o)[0] == qType) break;

        value+=file[o];
        escaped = false;
      }

      i++;
      curExp+=value;
      line+=value.match(/\n/g) == null ? 0 : value.match(/\n/g).length;

      lex.push({
        Name: '\'' + value + '\'',
        Exp: curExp,
        Line: line,
        Type: 'string',
        OName: '\'' + value + '\'',
        Dir: dir
      });
      continue outer;
    } else if (/^(\d|\+|\-|\.)/.test(file.substr(i))) {

      var sign = true;

      //detect positive and negative
      while (file.substr(i)[0] == '+' || file.substr(i)[0] == '-')
        if (file.substr(i) == '+') i++;
        else {
          sign = !sign;
          i++;
        }

      var num = '';

      if (!(/^(\d|\.)/.test(file.substr(i)))) {
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

      while (/^(\d|\.)/.test(file.substr(i))) {
        num+=file[i];
        i++;

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

      var variable = '$';

      var_loop:
        for (let o = i; o < file.length; o++) {

          for (let j = 0; j < KEYWORDS.length; j++)
            if (testkey(KEYWORDS[j], file, o)) break var_loop;

          variable+=file[o];
          i++;
        }

      curExp+=variable.substr(1);

      i--;
      lex.push({
        Name: variable,
        Exp: curExp,
        Line: line,
        Type: 'Variable',
        OName: variable.substr(1),
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
        errors.push(
          `Error while lexing in ${dir}! \nUnexpected token: \"${v.OName.trim()}\" on line ${v.Line}`
          +
          `\n\nFound near: ${v.Exp.trim()}`
          +
          `\n            ${' '.repeat(v.Exp.indexOf(v.OName) == -1 ? '' : v.Exp.indexOf(v.OName))}${'^'.repeat(v.OName.length)}`
          +
          `\nAdvanced Info: Two tokens of the type \"${v.Type}\" were detected adjacent to each other`
        );
      }
    }

    if (lex[i].Name.startsWith('$return'))
      //give a warning
      warnings.push(`Warning while lexing in ${dir}! Did you mean \"return ~\" on line ${lex[i].Line}?`);
  });

  return include_parser(dir, processes.insert_hashes(lex));
}

var stdinBuffer = fs.readFileSync(0)
, { f, dir, name } = JSON.parse(stdinBuffer.toString());

var
  warnings = [],
  errors = [];

console.log(
  JSON.stringify({
    WARNS: warnings,
    ERRORS: errors,
    LEX: lexer(processes.init(rw(f)), dir + name)
  }, null, 2)
);
