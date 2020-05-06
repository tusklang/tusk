#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
using namespace std;

Returner parser(const json actions, const json cli_params, json vars, const bool groupReturn, const bool expReturn);

json indexesCalc(json val, json indexes, json cli_params, json vars) {

  if (indexes.size() == 0) return val;

  string index = parser(indexes[0], cli_params, vars, false, true).exp["ExpStr"][0].get<string>();

  if (val.find(index) == val.end()) {
    if (val.find("falsey") == val.end()) return falseyVal;
    else return val["falsey"][0];
  } else {

    json expVal = parser(val[index], cli_params, vars, false, true).exp;

    indexes.erase(indexes.begin());

    if (indexes.size() == 0) return expVal;

    return indexesCalc(expVal["Hash_Values"], indexes, cli_params, vars);
  }
}

#endif
