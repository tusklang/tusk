//file to insert newlineS' (aka expression terminators)

module.exports = lex => {

  var newLex = [];

  for (let i = 0; i < lex.length; i++) {

    let currentType = lex[i].Type == 'operation' || lex[i].Type == '?operation' || lex[i].Type == 'namespace';
    let nextType = lex[i + 1] && (lex[i + 1].Type == 'operation' || lex[i + 1].Type == '?operation' || lex[i + 1].Type == 'namespace');

    newLex.push(lex[i]);

    if (lex[i].Name == '[:' || lex[i].Name == '(' || lex[i].Name == '[' || lex[i].Name == '{') continue; //because opening braces dont need semicolons after them, but if they do it is an error

    //because if and functions needs a ( after if
    if ((lex[i].Type == 'cond' || lex[i].Name == 'function') && lex[i + 1] && lex[i + 1].Name == '(') continue;

    //detect a type with the ? prefix
    if (lex[i].Type.startsWith('?') && (lex[i + 1] && lex[i + 1].Type.startsWith('?'))) continue;

    //example
    //  a['f' <-- expression value]
    //  would insert a newlineS
    //also
    //  abc: 'ff'
    //  log abc
    //would not insert a newlineS
    if (
      (lex[i].Type == 'expression value' && lex[i + 1] && lex[i + 1].Type == '?operation_paren')
      ||
      (lex[i].Type == '?operation_paren' && lex[i + 1] && lex[i + 1].Type == 'expression value')

      //but if it looks like:
      //  abc: [:
      //    a: 'a'
      //  :]
      //  abcd: 'ss'
      //it must insert a newlineS

      &&

      (function() {

        for (let o = i; o < lex.length; ++o) {

          if (lex[o].Name == '{') return true;
          if (lex[o].Name == '}') break;

          if (lex[o].Name == '[:') return true;
          if (lex[o].Name == ':]') break;

          if (lex[o].Name == '[') return true;
          if (lex[o].Name == ']') break;

          if (lex[o].Name == '(') return true;
          if (lex[o].Name == ')') break;

          if (cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && (lex[o].Type == 'operation' || lex[o].Type == '?operation' || lex[o].Type == 'namespace' || lex[o].Type == 'expression value') && lex[o].Name != '.') return false;

          if (lex[o].Name == ':') return true;
        }

        return false;

      }())

    ) continue;

    if (currentType == nextType) newLex.push({
      Name: 'newlineS',
      Exp: lex[i].Exp,
      Line: lex[i].Line,
      Type: 'operation',
      OName: ';',
      Dir: lex[i].Dir
    });
  }

  return newLex;
}
