module.exports = (key, file, i) => {
  if (file.substr(i).match(new RegExp('^' + key.pattern))) {
    return true;
  }

  return false;
}
