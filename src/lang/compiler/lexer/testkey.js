module.exports = (key, file, i) => {

  var matches = [];

  var re = new RegExp(key.pattern, 'g');

  while ((match = re.exec(file)) != null) matches.push(match.index);

  if (matches.includes(i)) {

    let pre_i = file.substr(0, i);

    if (key.no_id_before) {

      var ids = require('./keywords.json').filter(k => k.type == 'id' || k.type == 'id_non_tilde');

      ids = ids.filter(i => !!pre_i.match(new RegExp('(' + i.pattern + ')$')));

      if (ids.length != 0) return false;
    }

    if (key.no_oper_before) {

      var opers = require('./keywords.json').filter(k => k.type == 'operation' || k.type == '?operation' || k.type == '?operation_paren');

      opers = opers.filter(i => !!pre_i.trim().match(new RegExp('(' + i.pattern + ')$')));

      if (opers.length != 0) return false;
    }

    return true;
  }

  return false;
}
