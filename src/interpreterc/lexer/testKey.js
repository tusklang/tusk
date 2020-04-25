function expFormat(input) {

  if (input == '$variable') return '(\\$.*)';
  if (input == '(') return '\\(';

  return input;
}

function format(input) {

  if (input == '\n') return 'newline';
  if (input == '\\d') return 'digit';
  if (input == '\\w') return 'word';

  if (input == '$variable') return 'variable';

  if (input.startsWith('\\')) input = input.substr(1);

  return input;
}

module.exports = (keeper, file, keyword, line, curExp) => {
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

          var first = keeper.slice(0, -1 * file.length);

          if (!first) return false;

        } else {

          var pattern = new RegExp(keyword.pre_prohib[i] + '$')
          , first = keeper.slice(0, -1 * file.length)

          //create an "array" of needed items
          , needed_items = '[';

          for (let i = 0; i < keyword.pre_prohib.length; i++) needed_items+=format(keyword.pre_prohib[i]) + (i + 1 == keyword.post_or_necc.length ? '' : ', ');

          needed_items+=']';

          if (keyword.pre_prohib.length == 1) needed_items = format(keyword.pre_prohib[0]);

          if (first.match(pattern)) {

            //give an error
            console.log(
              `Error During Lexing >> Expected${keyword.pre_prohib.length != 1 ? ' One Of The Following' : ''}:`,
              needed_items,
              `${keyword.pre_prohib.length != 1 ? '\n' : ''}After \"${keyword.remove}\"`,
              `\n\nError Occurred On Line: ${line}`
            );

            process.exit(1);

            return false;
          }
        }
      }
    }

    if (keyword.pre_or_necc) {

      var or_exp = '(';

      for (let i = 0; i < keyword.pre_or_necc.length; i++) or_exp+=expFormat(keyword.pre_or_necc[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      var pattern = new RegExp(or_exp + '$')
      , first = keeper.slice(0, -1 * file.length)
      , needed_items = '[';

      for (let i = 0; i < keyword.pre_or_necc.length; i++) needed_items+=format(keyword.pre_or_necc[i]) + (i + 1 == keyword.pre_or_necc.length ? '' : ', ');

      needed_items+=']';

      if (keyword.pre_or_necc.length == 1) needed_items = format(keyword.pre_or_necc[0]);

      if (!first.match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected${keyword.pre_or_necc.length != 1 ? ' One Of The Following' : ''}:`,
          needed_items,
          `${keyword.pre_or_necc.length != 1 ? '\n' : ''}Before \"${keyword.remove}\"`,
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc) {

      //create the or exp
      var or_exp = '(';

      for (let i = 0; i < keyword.post_or_necc.length; i++) or_exp+=expFormat(keyword.post_or_necc[i]) + '|';

      or_exp = or_exp.slice(0, -1);
      or_exp+=')';

      //create the pattern
      var pattern = new RegExp('^' + or_exp)

      //create an "array" of needed items
      , needed_items = '[';

      for (let i = 0; i < keyword.post_or_necc.length; i++) needed_items+=format(keyword.post_or_necc[i]) + (i + 1 == keyword.post_or_necc.length ? '' : ', ');

      needed_items+=']';

      if (keyword.post_or_necc.length == 1) needed_items = format(keyword.post_or_necc[0]);

      if (!file.substr(keyword.remove.length).match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected${keyword.post_or_necc.length != 1 ? ' One Of The Following' : ''}:`,
          needed_items,
          `${keyword.post_or_necc.length != 1 ? '\n' : ''}Before Expression: \"${file.substr(keyword.remove.length).slice(0, MAX_CUR_EXP_SIZE)}\"`,
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc_after_tilde) {

      file = file.substr(1);

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

      if (!file.substr(keyword.remove.length).match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected ${needed_items} After ${keyword.remove + '~'}`,
          `\n\nError Occurred On Line: ${line}`
        );

        process.exit(1);
        return false;
      }
    }

    if (keyword.post_or_necc_after_paren) {

      file = file.substr(1);

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

      if (!file.substr(keyword.remove.length).match(pattern)) {

        //give an error
        console.log(
          `Error During Lexing >> Expected${needed_items.startsWith('[') ? ' One Of The Following: \n' : ''}${needed_items}${needed_items.startsWith('[') ? '\n' : ' '}After ${keyword.remove + '('}`,
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
