#ifndef ARRAY_HPP_
#define ARRAY_HPP_

#include "../json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json arrayCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp[0][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2);

  //OPTIMIZE:

  while (IsLessC("0", &index[0])) {
    val.erase(val.begin());
    index = Subtract(&index[0], "1", &calc_params.dump()[0], line);
  }

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return parser(val[0], calc_params, vars, dir, false, line, true).exp[0][0];

  return indexesCalc(val[0][0]["Index_Type"] == "hash" ? val[0][0]["Hash_Values"] : val[0][0]["Value"], indexes, calc_params, vars, line, dir, val[0][0]["Index_Type"]);
}

#endif
