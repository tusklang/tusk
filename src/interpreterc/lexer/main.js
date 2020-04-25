const keywords = require('./keywords.json')
, testKey = require('./testKey.js')
, fs = require('fs');

//set the maximum size of the cur exp
global.MAX_CUR_EXP_SIZE = 20;

//fetch file
var stdinBuffer = fs.readFileSync(0)
, file = stdinBuffer.toString();

file = require('./procInit')(file);

const copyFile = file;

var lex = []
, line = 1

//cur exp is the current expression (to help locate the error)
, curExp = '';

for (let i = 0; i < keywords.length; i++) {

  if (file.length == 0) break;

  while (curExp.length > MAX_CUR_EXP_SIZE) curExp = curExp.substr(1);

  if (testKey(copyFile, file, keywords[i], line, curExp)) {

    curExp+=keywords[i].name;

    //if it is a newline then increment line
    if (keywords[i].name == 'newlineN') line++;

    file = file.substr(keywords[i].remove.length);
    lex.push({
      Name: keywords[i].name,
      Exp: curExp,
      Line: line
    });
    i = -1;
    continue;
  }

  if (i + 1 == keywords.length) {

    if (/^(\'|\"|\`)/.test(file)) {

      var typeOfQ = ''
      , exp = ''

      for (let o = 0; o < file.length; o++) {
        if (file[o] == '\'') {
          if (typeOfQ == '\'') {
            typeOfQ = '';
          } else if (typeOfQ == '') {
            typeOfQ = '\'';
          }
        }

        if (file[o] == '\"') {
          if (typeOfQ == '\"') {
            typeOfQ = '';
          } else if (typeOfQ == '') {
            typeOfQ = '\"';
          }
        }

        if (file[o] == '\`') {
          if (typeOfQ == '\`') {
            typeOfQ = '';
          } else if (typeOfQ == '') {
            typeOfQ = '\`';
          }
        }

        exp+=file[o];

        if (typeOfQ == '') {
          break;
        }
      }

      //increase line by amount of newlines in the string
      line+=exp.split('\n').length - 1;

      curExp+=exp;

      lex.push({
        Name: '\'' + exp.substr(1).slice(0, -1) + '\'',
        Exp: curExp,
        Line: line
      });

      file = file.substr(exp.length);

      i = -1;
      continue;
    } else if (/^(\d|\-|\+|\.)/.test(file)) {

      var num = '';

      if (/^(\-|\+)/.test(file)) {
        num+=file[0];
        file = file.substr(1)
      }

      while (/^(\d|\.)+/.test(file)) {
        num+=file[0];

        file = file.substr(1);
      }

      if (num == '-' || num == '+') num+='1';

      if (num.endsWith('.') && file.startsWith('[')) {
        num = num.slice(0, -1);
        lex.push({
          Name: num,
          Exp: curExp,
          Line: line
        });
        lex.push({
          Name: '.',
          Exp: curExp,
          Line: line
        });
      } else lex.push({
        Name: num,
        Exp: curExp,
        Line: line
      });

      curExp+=num;

      if (file.startsWith('(')) lex.push({
        Name: '*',
        Exp: curExp,
        Line: line
      });

      i = -1;
      continue;
    } else {

      var name = '';

      count: for (let o = 0; o < file.length; o++) {

        for (let j = 0; j < keywords.length; j++) {
          let pattern = new RegExp(keywords[j].pattern);

          if (file.substr(o).search(pattern) == 0) break count;
        }

        name+=file[o];
      }

      file = file.substr(name.length);
      lex.push({
        Name: '$' + name,
        Exp: curExp,
        Line: line
      });

      curExp+='$' + name;

      i = -1;
      continue;
    }
  }
}

// send data through stdout back to go process
console.log(JSON.stringify(lex));
