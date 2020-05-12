// this will allow you to write #procName instead of #~procName
module.exports.init = file => {

  var nFile = ""
  , escaped = false
  , typeOfQ = "";

  for (let i = 0; i < file.length; i++) {
    if (!escaped && file[i] == '\\') {
      escaped = true;
      i--;
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

//this will allow you to write procName() instead of #procName() or #~procName()
module.exports.insert_hashes = lex => {

  for (let i = 0; i < lex.length; i++)
    if (lex[i].Name.startsWith('$') && lex[i + 1] && lex[i + 1].Name == '(') {

      if (lex[i - 2] && (lex[i - 2].Name == '#' || lex[i - 2].Name == '@' || lex[i - 2].Name == 'process')) continue;

      if (i - 1 < 0) lex.unshift({
        Name: '#',
        Exp: '#~' + lex[i].Exp,
        Line: lex[i].Line,
        Type: 'id',
        OName: '#',
        Dir: lex[i].Dir
      }, {
        Name: '~',
        Exp: '~' + lex[i].Exp,
        Line: lex[i].Line,
        Type: 'operation',
        OName: '~',
        Dir: lex[i].Dir
      }); else lex.splice(i, 0, {
        Name: '#',
        Exp: '#~' + lex[i].Exp,
        Line: lex[i].Line,
        Type: 'id',
        OName: '#',
        Dir: lex[i].Dir
      }, {
        Name: '~',
        Exp: '~' + lex[i].Exp,
        Line: lex[i].Line,
        Type: 'operation',
        OName: '~',
        Dir: lex[i].Dir
      });

      i+=2;
    } else if (lex[i].Name == ']' && lex[i + 1] && lex[i + 1].Name == '(') {

      var o;

      var bCnt = 0;

      for (o = i; o >= 0; o--) {

        if (lex[o].Name == '[') bCnt++;
        if (lex[o].Name == ']') bCnt--;

        if (bCnt == 0 && lex[o - 1].Name == '.' && lex[o].Name == '[') break;
      }

      if (lex[o - 3] && lex[o - 3].Name == "~" && lex[o - 4].Name && lex[o - 4].Name == "#") continue;

      lex.splice(o - 2, 0, {
        Name: '#',
        Exp: '#~' + lex[o].Exp,
        Line: lex[o].Line,
        Type: 'id',
        OName: '#',
        Dir: lex[o].Dir
      }, {
        Name: '~',
        Exp: '~' + lex[o].Exp,
        Line: lex[o].Line,
        Type: 'operation',
        OName: '~',
        Dir: lex[o].Dir
      });

      i+=2;
    }

  return lex;
}
