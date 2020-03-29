#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include <iostream>
#include "json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, int line, string type);
json arrayCalc(json val, json indexes, json calc_params, int line);

json arrayCalc(json val, json indexes, json calc_params, int line) {
  string index = indexes[0].dump().substr(1, indexes[0].dump().length() - 2);

  while (IsLessC("0", &index[0])) val.erase(val.begin());

  indexes.erase(indexes.begin());

  return indexesCalc(val[0]["Value"], indexes, calc_params, line, val[0]["Index_Type"]);
}

json indexesCalc(json val, json indexes, json calc_params, int line, string type) {

  if (indexes.size() == 0) return val;

  if (type == "array") return arrayCalc(val, indexes, calc_params, line);

  return "{}"_json;
}

#endif
