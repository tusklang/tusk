const keywords = require('./keywords.json');

//this entire file is bad and needs re-factoring

function expFormat(input) {

  if (input == '$variable') {

    var getKeyWords = keywords.map(k => {

      if (k == '\r\n' || k == '\n') return '\\n';
      else return k.pattern;
    }).join('|')
    , matcher = `(?!(${getKeyWords}))`;

    return matcher;
  }

  if (input == '$none') return '\\s*$';
  if (input == '$notnone') return '.';
  if (input == '$string') return '((\'|\"|\`)[\\s\\S]+(\'|\"|\`))';
  if (input == '$hash') return '((\\[\\:)[\\s\\S]+(\\:\\]))';
  if (input == '$array') return '((\\[)[\\s\\S]+(\\]))';
  if (input == '(') return '\\(';
  if (input == ')') return '\\)';
  if (input == ']') return '\\]';
  if (input == ':]') return '\\:\\]'
  if (input == '*keyword') {

    //craft the regex of all the keywords that are ids
    var ids = keywords.filter(k => k.type == 'id').map(k => "(" + k.remove + ")").join('|')
    , match_reg = `(${ids})`;

    return match_reg;
  }
  if (input == '*operation') {

    var opers = keywords.filter(k => k.type == 'operation').map(k => "(" + k.pattern + ")").join('|')
    , match_reg = `(${opers})`;

    return match_reg;
  }
  if (input == '*comparison') {

    var comps = keywords.filter(k => k.type == 'comparison').map(k => "(" + k.pattern + ")").join('|')
    , match_reg = `(${comps})`;

    return match_reg;
  }
  if (input == '*math') {

    var math = keywords.filter(k => k.type == 'math').map(k => "(" + k.pattern + ")").join('|')
    , match_reg = `(${math})`;

    return match_reg;
  }

  return input;
}

function format(input) {

  if (input == '\n') return 'newline';
  if (input == '\\d') return 'digit';
  if (input == '\\w') return 'word';

  if (input == '$none') return 'none';
  if (input == '$notnone') return 'not none';
  if (input == '$variable') return 'variable';
  if (input == '$string') return 'string';
  if (input == '$hash') return 'hash';
  if (input == '$array') return 'array';
  if (input == '*keyword') return 'Any Keyword';
  if (input == '*operation') return 'Any Operation';
  if (input == '*comparison') return 'Any Comparison';
  if (input == '*math') return 'Any Math Operation';

  if (input.includes('\\')) input = input.replace('\\', '');

  return input;
}

module.exports = (keeper, file, keyword, line, curExp, lex) => {
  var pattern = new RegExp('^' + keyword.pattern);

  if (file.match(pattern)) {

    if (keyword.pre_prohib_soft) {

      for (let i = 0; i < keyword.pre_prohib_soft.length; i++) {

        if (keyword.pre_prohib_soft[i] == "$none") {

          var first = keeper.slice(0, -1 * file.length);

          if (!first) return false;

        } else {

          var pattern = new RegExp(keyword.pre_prohib_soft[i] + '$')
          , first = keeper.slice(0, -1 * file.length);

          if (first.match(pattern)) {
            return false;
          }
        }
      }
    }

    //throw an error if the pre prohib is violated
    if (keyword.pre_prohib) {

      for (let i = 0; i < keyword.pre_prohib.length; i++) {

        if (keyword.pre_prohib[i] == "$none") {

          if (!lex[lex.length - 1]) return false;

        } else {

          if (!lex[lex.length - 1]) return true;

          var pattern = new RegExp(expFormat(keyword.pre_prohib[i]) + '$')

          //create an "array" of needed items
          , needed_items = '[';

          for (let i = 0; i < keyword.pre_prohib.length; i++) needed_items+=format(keyword.pre_prohib[i]) + (i + 1 == keyword.pre_prohib.length ? '' : ', ');

          needed_items+=']';

          if (keyword.pre_prohib.length == 1) needed_items = format(keyword.pre_prohib[0]);

          if (lex[lex.length - 1].Name.match(pattern)) {

            //give an error
            console.log(
              `Error During Lexing >> Expected${keyword.pre_prohib.length != 1 ? ' One Of The Following' : ''}:`,
              needed_items,
              `${keyword.pre_prohib.length != 1 ? '\n' : ''}After \"${keyword.remove}\"`,
              `\nInstead Got ${!keeper[keyword.remove.length + 1] ? 'nothing' : format(keeper[keyword.remove.length + 1])}`,
              `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
              `\n            ${' '.repeat(keyword.remove.length - 1) + '^~~~'}`, //point to error location
              `\n\nError Occurred On Line: ${line}`
            );

            process.exit(1);

            return false;
          }
        }
      }
    }

    if (keyword.pre_or_necc) {

      if (!lex[lex.length - 1]) {
        if (keyword.pre_or_necc.includes('$none')) console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        ); else return true;
      }

      var or_exp = '(';

      for (let i = 0; i < keyword.pre_or_necc.length; i++) or_exp+=expFormat(keyword.pre_or_necc[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      var pattern = new RegExp(or_exp + '$')
      , needed_items = '[';

      for (let i = 0; i < keyword.pre_or_necc.length; i++) needed_items+=format(keyword.pre_or_necc[i]) + (i + 1 == keyword.pre_or_necc.length ? '' : ', ');

      needed_items+=']';

      if (keyword.pre_or_necc.length == 1) needed_items = format(keyword.pre_or_necc[0]);

      if (!lex[lex.length - 1].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected${keyword.pre_or_necc.length != 1 ? ' One Of The Following' : ''}:`,
          needed_items,
          `${keyword.pre_or_necc.length != 1 ? '\n' : ''}Before \"${keyword.remove}\"`,
          `\nInstead Got ${curExp == '' ? 'nothing' : curExp}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${' '.repeat(keyword.remove.length - 1) + '^~~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc) {

      if (!lexer(file_)[0]) {
        if (keyword.post_prohib.includes('$none')) console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        ); else return true;
      }

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_or_necc.length; i++) or_exp+=expFormat(keyword.post_or_necc[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)
      , f_lex = lexer(file)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_or_necc.length; i++) needed_items+=format(keyword.post_or_necc[i]) + (i + 1 == keyword.post_or_necc.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_or_necc.length == 1) needed_items = format(keyword.post_or_necc[0]);

      if (!f_lex[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected${keyword.post_or_necc.length != 1 ? ' One Of The Following:' : ''}`,
          needed_items,
          `${keyword.post_or_necc.length != 1 ? '\n' : ''}After \"${keyword.remove}\"`,
          `\nInstead Got ${!file.slice(keyword.remove.length, MAX_CUR_EXP_SIZE) ? 'nothing' : file.slice(keyword.remove.length, MAX_CUR_EXP_SIZE)}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n            ${' '.repeat(keyword.remove.length - 1) + '^~~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc_after_tilde) {

      var file_ = file.substr(1);

      if (!lexer(file_)[0]) {
        if (keyword.post_prohib.includes('$none')) console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        ); else return true;
      }

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_or_necc_after_tilde.length; i++) or_exp+=expFormat(keyword.post_or_necc_after_tilde[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_or_necc_after_tilde.length; i++) needed_items+=format(keyword.post_or_necc_after_tilde[i]) + (i + 1 == keyword.post_or_necc_after_tilde.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_or_necc_after_tilde.length == 1) needed_items = format(keyword.post_or_necc_after_tilde[0]);

      if (!lexer(file_)[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove + '~'}`,
          `\nInstead Got ${!file_[keyword.remove.length + 1] ? 'nothing' : format(file_[keyword.remove.length + 1])}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n            ${' '.repeat(keyword.remove.length - 2) + '^~~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc_after_paren) {

      var file_ = file.substr(1);

      if (!lexer(file_)[0]) {
        if (keyword.post_prohib.includes('$none')) console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        ); else return true;
      }

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_or_necc_after_paren.length; i++) or_exp+=expFormat(keyword.post_or_necc_after_paren[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_or_necc_after_paren.length; i++) needed_items+=format(keyword.post_or_necc_after_paren[i]) + (i + 1 == keyword.post_or_necc_after_paren.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_or_necc_after_paren.length == 1) needed_items = format(keyword.post_or_necc_after_paren[0]);

      if (!lexer(file_)[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove + '('}`,
          `\nInstead Got ${!file_[keyword.remove.length + 1] ? 'nothing' : format(file_[keyword.remove.length + 1])}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '~~~^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_prohib) {

      var file_ = file;

      if (!lexer(file_.substr(keyword.remove.length))[0]) {
        if (keyword.post_prohib.includes('$none')) console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        ); else return true;
      }

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_prohib.length; i++) or_exp+=expFormat(keyword.post_prohib[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_prohib.length; i++) needed_items+=format(keyword.post_prohib[i]) + (i + 1 == keyword.post_prohib.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_prohib.length == 1) needed_items = format(keyword.post_prohib[0]);

      if (lexer(file_.substr(keyword.remove.length))[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_prohib_after_tilda) {

      var file_ = file.substr(1);

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_prohib_after_tilda.length; i++) or_exp+=expFormat(keyword.post_prohib_after_tilda[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_prohib_after_tilda.length; i++) needed_items+=format(keyword.post_prohib_after_tilda[i]) + (i + 1 == keyword.post_prohib_after_tilda.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_prohib_after_tilda.length == 1) needed_items = format(keyword.post_prohib_after_tilda[0]);

      if (lexer(file_)[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove + '~'}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '~~~^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_prohib_after_paren) {

      var file_ = file.substr(1);

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_prohib_after_paren.length; i++) or_exp+=expFormat(keyword.post_prohib_after_paren[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_prohib_after_paren.length; i++) needed_items+=format(keyword.post_prohib_after_paren[i]) + (i + 1 == keyword.post_prohib_after_paren.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_prohib_after_paren.length == 1) needed_items = format(keyword.post_prohib_after_paren[0]);

      if (lexer(file_)[0].Name.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Unexpected ${needed_items.startsWith('[') ? 'One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove + '('}`,
          `\nInstead Got ${!file_[keyword.remove.length + 1] ? 'nothing' : format(file_[keyword.remove.length + 1])}`,
          `\n\nFound Near: \"${file.slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n             ${'~'.repeat(keyword.remove.length - 2) + '~~~^~~'}`, //point to error location
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    return true;
  }

  return false;
}
