module.exports = lex => {

  var nLex = [];

  for (let i = 0; i < lex.length; i++) {

    nLex.push(lex[i]);

    if (lex[i].Type == 'id' && lex[i + 1] && lex[i + 1].Name != '~') nLex.push({
      Name: '~',
      Exp: lex[i].Exp,
      Line: lex[i].Line,
      Type: 'operation',
      OName: '~',
      Dir: lex[i].Dir
    });
  }

  return nLex;
}
