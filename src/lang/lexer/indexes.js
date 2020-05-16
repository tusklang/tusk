module.exports = (lex) => {

  var index;

  for (let i = 0; i < lex.length; i++) {

    if (lex[i + 1] && lex[i].Name == '::') index = [{
      Name: '.',
      Exp: lex[i].Exp,
      Line: lex[i].line,
      Type: 'operation',
      OName: '.',
      Dir: lex[i].Dir
    }, {
      Name: '[',
      Exp: lex[i].Exp,
      Line: lex[i].line,
      Type: '?operation',
      OName: '[',
      Dir: lex[i].Dir
    }, {
      Name: `\'${lex[i + 1].Name.substr(1)}\'`,
      Exp: lex[i + 1].Exp,
      Line: lex[i + 1].line,
      Type: 'string',
      OName: `\'${lex[i + 1].Name.substr(1)}\'`,
      Dir: lex[i + 1].Dir
    }, {
      Name: ']',
      Exp: lex[i + 1].Exp,
      Line: lex[i + 1].line,
      Type: '?operation',
      OName: ']',
      Dir: lex[i + 1].Dir
    }];

    if (index) lex.splice(i, 2, ...index);

    index = undefined;
  }

  return  lex;
}
