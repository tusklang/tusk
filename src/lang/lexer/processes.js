const keywords = require('./keywords.json');

// this will allow you to write #procName instead of #~procName
module.exports.init = file => {

  var nFile = ""
  , escaped = false
  , typeOfQ = "";

  for (let i = 0; i < file.length; i++) {
    if (!escaped && file[i] == '\\') {
      escaped = true;
      continue;
    }

    if (!escaped) {
      if (/['`"]/.test(file[i])) {

        if (typeOfQ == '') typeOfQ = file[i];
        else if (typeOfQ == file[i]) typeOfQ = '';
      }
    } else escaped = false;

    if (typeOfQ == '' && (file[i] == '#' || file[i] == '@') && file[i + 1] != '~') nFile+=file[i] + '~';
    else nFile+=file[i];
  }

  return nFile;
}

module.exports.hash_inserter = lex => {

  var cprocs = keywords.filter(k => k.type == "cproc").map(k => k.name);

  for (let i = 0; i < lex.length; i++) {

    if (lex[i].Name.startsWith('$')) {

      //make sure it is not a process name
      if (lex[i - 2] && lex[i - 2].Name == "process") continue;

      var insert_hash = false;

      var
        glCnt = 0,
        cbCnt = 0,
        bCnt = 0,
        pCnt = 0;

      for (let o = i + 1; o < lex.length; o++) {

        if (lex[o].Name == '[:') glCnt++;
        if (lex[o].Name == ':]') glCnt--;

        if (lex[o].Name == '{') cbCnt++;
        if (lex[o].Name == '}') cbCnt--;

        if (lex[o].Name == '[') bCnt++;
        if (lex[o].Name == ']') bCnt--;

        if (glCnt != 0 || cbCnt != 0 || bCnt != 0 || pCnt != 0) continue;

        if (lex[o].Name == '(') {
          insert_hash = true;
          break;
        }

        if (lex[o].Name != '.' && lex[o].Name != '[' && lex[o].Name != ']') break;
      }

      if (insert_hash) {

        let inserter = [
          {
            Name: '#',
            Exp: lex[i].Exp,
            Line: lex[i].Line,
            Type: 'caller',
            OName: '#',
            Dir: lex[i].Dir
          },
          {
            Name: '~',
            Exp: lex[i].Exp,
            Line: lex[i].Line,
            Type: 'operation',
            OName: '~',
            Dir: lex[i].Dir
          }
        ];

        if (i == 0) lex.unshift(...inserter)
        else lex.splice(i, 0, ...inserter);

        i+=2
      }

    }
  }

  return lex;
}
