#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include "json.hpp"
#include "parser.hpp"
using namespace std;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn);

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  if (indexes.size() == 0) return val;

  string index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

  if (val.find(index) == val.end()) return val["falsey"][0];
  else {

    json expVal = parser(val[index], calc_params, vars, dir, false, line, true).exp;

    indexes.erase(indexes.begin());

    if (indexes.size() == 0) return expVal;

    return indexesCalc(expVal["Hash_Values"], indexes, calc_params, vars, line, dir);
  }
}

#endif
