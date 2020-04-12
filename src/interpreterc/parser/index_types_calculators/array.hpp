#ifndef ARRAY_HPP_
#define ARRAY_HPP_

#include "../parser.hpp"
#include "../json.hpp"
#include "../structs.h"
#include "../indexes.hpp"
using namespace std;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn);
json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json arrayCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2);

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return val[index][0];

  return indexesCalc(val[index][0]["Hash_Values"], indexes, calc_params, vars, line, dir, val[index][0]["Index_Type"]);
}

#endif
