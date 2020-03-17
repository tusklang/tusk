const keywords = require('./keywords.json')
, testKey = require('./testKey.js')
, fs = require('fs');

//fetch file
var stdinBuffer = fs.readFileSync(0)
, file = stdinBuffer.toString();

file = require('./procInit')(file);

const copyFile = file;

var lex = [];

for (let i = 0; i < keywords.length; i++) {

  if (file.length == 0) break;

  if (testKey(copyFile, file, keywords[i])) {

    file = file.substr(keywords[i].remove.length);
    lex.push(keywords[i].name);
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

      lex.push('\'' + exp.substr(1).slice(0, -1) + '\'');

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

      if (num == "-" || num == "+") num+="1";

      lex.push(num);

      if (file.startsWith("(")) lex.push("*");

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
      lex.push('$' + name);

      i = -1;
      continue;
    }
  }
}

// send data through stdout back to go process
console.log(JSON.stringify(lex));
