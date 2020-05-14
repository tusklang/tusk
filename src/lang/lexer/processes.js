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
