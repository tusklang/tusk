module.exports = lex => {

  var index;

  for (let i = 0; i < lex.length; i++) {

    if (lex[i + 1] && lex[i].Name == '::') index = [{
      Name: '.',
      Exp: lex[i].Exp,
      Line: lex[i].Line,
      Type: 'operation',
      OName: '.',
      Dir: lex[i].Dir
    }, {
      Name: '[',
      Exp: lex[i].Exp,
      Line: lex[i].Line,
      Type: '?operation',
      OName: '[',
      Dir: lex[i].Dir
    }, {
      Name: `\'${lex[i].Name[0] == '$' ? lex[i + 1].Name.substr(1) : lex[i + 1].Name}\'`,
      Exp: lex[i + 1].Exp,
      Line: lex[i + 1].Line,
      Type: 'string',
      OName: `\'${lex[i].Name[0] == '$' ? lex[i + 1].Name.substr(1) : lex[i + 1].Name}\'`,
      Dir: lex[i + 1].Dir
    }, {
      Name: ']',
      Exp: lex[i + 1].Exp,
      Line: lex[i + 1].Line,
      Type: '?operation',
      OName: ']',
      Dir: lex[i + 1].Dir
    }];

    if (index) lex.splice(i, 2, ...index);

    index = undefined;
  }

  return  lex;
}
