module.exports = (key, file, i) => {

  var matches = [];

  var re = new RegExp(key.pattern, 'g');

  while ((match = re.exec(file)) != null) matches.push(match.index);

  if (matches.includes(i)) return true;

  return false;
}
