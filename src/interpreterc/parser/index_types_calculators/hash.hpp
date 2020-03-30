#ifndef HASH_HPP_
#define HASH_HPP_

#include "../json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json hashCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp[0][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2);

  json indexVal = val[index];

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return parser(indexVal, calc_params, vars, dir, false, line, true).exp[0][0];

  return indexesCalc(indexVal[0]["Index_Type"] == "hash" ? indexVal[0]["Hash_Values"] : indexVal[0]["Value"], indexes, calc_params, vars, line, dir, indexVal[0]["Index_Type"]);
}

#endif
