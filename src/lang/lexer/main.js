const
  testkey = require('./testkey'),
  processes = require('./processes')
  fs = require('fs');

global.KEYWORDS = require('./keywords.json');
global.MAX_CUR_EXP = 12;

var lexer = (file) => {

  //current expression
  var curExp = ''
  //current lexxed value
  , lex = []
  //current line
  , line = 1;

  //loop through file
  outer:
  for (let i = 0; i < file.length; i++) {

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
          OName: KEYWORDS[o].remove
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
        if (escaped) escaped = false;

        //if it is escape character, set escaped to true
        if (file[o] == '\\') escaped = true;

        if (file.substr(o)[0] == qType && !escaped) break;

        value+=file[o];
      }

      i++;
      curExp+=value;
      line+=value.match(/\n/g) == null ? 0 : value.match(/\n/g).length;

      lex.push({
        Name: '\'' + value + '\'',
        Exp: curExp,
        Line: line,
        Type: 'string',
        OName: '\'' + value + '\''
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
          OName: num
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
        OName: num
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
        OName: variable.substr(1)
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
        console.log(
          `Error while lexing! \nUnexpected token: \"${v.OName.trim()}\" on line ${v.Line}`,
          `\n\nFound near: ${v.Exp.trim()}`,
          `\n            ${' '.repeat(v.Exp.indexOf(v.OName) == -1 ? '' : v.Exp.indexOf(v.OName))}${'^'.repeat(v.OName.length)}`,
          `\nAdvanced Info: Two tokens of the type \"${v.Type}\" were detected adjacent to each other`
        );
        process.exit(1);
      }
    }
  });

  return lex;
}

var stdinBuffer = fs.readFileSync(0)
, f = stdinBuffer.toString();

console.log(
  JSON.stringify(
    processes.insert_hashes(
      lexer(
        processes.init(f)
      )
    )
  )
);
