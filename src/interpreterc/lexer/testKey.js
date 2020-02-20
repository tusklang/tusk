module.exports = (keeper, file, keyword) => {
  var pattern = new RegExp('^' + keyword.pattern);

  if (file.match(pattern)) {

    if (keyword.pre_prohib_soft) {
      var pattern = new RegExp(keyword.pre_prohib_soft + '$')
      , first = keeper.slice(0, -1 * file.length)

      if (first.match(pattern)) {
        return false
      }
    }

    return true;
  }

  return false;
}
