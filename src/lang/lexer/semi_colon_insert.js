module.exports = lex => {
  
  var newLex = [];

  for (let i = 0; i < lex.length; i++) {

    let currentType = lex[i].Type == 'operation' || lex[i] == '?operation' || lex[i].Type == 'namespace';
    let nextType = lex[i + 1] && (lex[i + 1].Type == 'operation' || lex[i + 1] == '?operation' || lex[i + 1].Type == 'namespace');

    newLex.push(lex[i]);

    if (lex[i].Type == 'expression value' && (lex[i + 1] && lex[i + 1].Type == '?operation')) newLex.push({
      Name: 'newlineS',
      Exp: lex[i].Exp,
      Line: lex[i].Line,
      Type: 'operation',
      OName: ';',
      Dir: lex[i].Dir
    });

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
