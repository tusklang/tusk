module.exports = lex => {

  var newLex = [];

  for (let i = 0; i < lex.length; i++) {

    let currentType = lex[i].Type == 'operation' || lex[i].Type == '?operation' || lex[i].Type == 'namespace' || lex[i].Type == '?paren' || lex[i].Type == '?operation_paren';
    let nextType = lex[i + 1] && (lex[i + 1].Type == 'operation' || lex[i + 1].Type == '?operation' || lex[i + 1].Type == 'namespace' || lex[i + 1].Type == '?paren' || lex[i + 1].Type == '?operation_paren');

    newLex.push(lex[i]);

    //detect a type with the ? prefix
    if (lex[i].Type.startsWith('?') && (lex[i + 1] && lex[i + 1].Type.startsWith('?'))) continue;

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
