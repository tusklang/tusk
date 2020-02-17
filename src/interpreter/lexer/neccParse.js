var operators = ['+', '-', '*', '/', '^', '%', '~', '>', '<', '=', '!', ':', ';']
, keywords = require('./keywords.json');

module.exports = (file, key) => {

  while (key.includes('*operators')) {
    key = key.splice(key.indexOf('*operators'), 1);

    if (operators.includes(files[0])) return true;
  }

  while (key.includes('*keywords')) {
    key = key.splice(key.indexOf('*keywords'), 1);

    let patterns = keywords.map(k => new RegExp(k.pattern).test(file));

    if (patterns.includes(true)) return true;
  }

  for (let i = 0; i < key.length; i++) {
    var pattern = new RegExp('^' + key[i]);

    if (pattern.test(file)) {
      return false;
    }
  }

  return false;
}
