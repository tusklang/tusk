#ifndef INDEXES_HPP_
#define INDEXES_HPP_

#include <map>
#include <vector>
#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

Returner parser(const vector<Action> actions, const json cli_params, map<string, Variable> vars, const bool groupReturn, const bool expReturn);

Action indexesCalc(map<string, vector<Action>> val, vector<vector<Action>> indexes, json cli_params, map<string, Variable> vars) {

  if (indexes.size() == 0) return falseyVal;

  string index = parser(indexes[0], cli_params, vars, false, true).exp.ExpStr[0];

  if (val.find(index) == val.end()) {

    if (val.find("falsey") == val.end()) return falseyVal;
    else return val["falsey"][0];

  } else {

    Action expVal = parser(val[index], cli_params, vars, false, true).exp;

    indexes.erase(indexes.begin());

    if (indexes.size() == 0) return expVal;

    return indexesCalc(expVal.Hash_Values, indexes, cli_params, vars);
  }
}

#endif
