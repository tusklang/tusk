module.exports = (keeper, file, keyword) => {
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
            return false
          }
        }
      }
    }

    return true;
  }

  return false;
}
