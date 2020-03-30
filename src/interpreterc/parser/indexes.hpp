#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include "json.hpp"

#include "index_types_calculators/array.hpp"
#include "index_types_calculators/hash.hpp"

using namespace std;

json expressionCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {
  return "{}"_json;
}

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type) {

  if (indexes.size() == 0) return val;

  if (type == "array") return arrayCalc(val, indexes, calc_params, vars, line, dir);
  else if (type == "hash") return hashCalc(val, indexes, calc_params, vars, line, dir);
}

#endif
