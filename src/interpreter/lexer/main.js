const keywords = require('./keywords.json')
, neccParse = require('./neccParse')
, fs = require('fs');

// var stdinBuffer = fs.readFileSync(0)
// , file = stdinBuffer.toString();

file = '~~~';

file = require('./procInit')(file);

var lex = [];

for (let i = 0; i < keywords.length; i++) {

  if (file.length == 0) break;

  var pattern = new RegExp(keywords[i].pattern);

  if (pattern.test(file)) {

    if (keywords[i].after_soft_necc && !neccParse(file.substr(keywords[i].remove.length), keywords[i].after_soft_necc)) continue;
    if (keywords[i].after_necc && !neccParse(file.substr(keywords[i].remove.length), keywords[i].after_necc)) {

      if (file.substr(keywords[i].remove.length)[0]) console.log('Error: Unexpected ' + '\'' + file.substr(keywords[i].remove.length)[0] + '\'');
      else console.log('Error: Expected Something After ' + keywords[i].remove);

      process.exit(1);
    } else if (keywords[i].after_prohib && neccParse(file.substr(keywords[i].remove.length), keywords[i].after_prohib)) {

      if (file.substr(keywords[i].remove.length)[0]) console.log('Error: Unexpected ' + '\'' + file.substr(keywords[i].remove.length)[0] + '\'');
      else console.log('Error: Expected Something After ' + keywords[i].remove);

      process.exit(1);
    } else {
      file = file.substr(keywords[i].remove.length);
      lex.push(keywords[i].name);
      i = -1;
      continue;
    }
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

      lex.push(exp);

      file = file.substr(exp.length);

      i = -1;
      continue;
    } else if (/^\d/.test(file)) {
      var num = '';

      while (/^\d+/.test(file)) {
        num+=file[0];

        file = file.substr(1);
      }

      lex.push(num);

      i = -1;
      continue;
    } else {

      var name = '';

      count: for (let o = 0; o < file.length; o++) {

        for (let j = 0; j < keywords.length; j++) {
          let pattern = new RegExp(keywords[j].pattern);

          if (pattern.test(file.substr(o))) break count;
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
