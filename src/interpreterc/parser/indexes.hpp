#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include "json.hpp"

#include <iostream>
#include "index_types_calculators/array.hpp"
#include "index_types_calculators/hash.hpp"
#include "index_types_calculators/expressions.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type) {

  if (indexes.size() == 0) return val;

  if (type == "array") return arrayCalc(val, indexes, calc_params, vars, line, dir);
  else if (type == "hash") return hashCalc(val, indexes, calc_params, vars, line, dir);
  else if (type == "expression") return expressionCalc(val["ExpAct"], indexes, calc_params, vars, line, dir);

  return "{}"_json;
}

#endif
