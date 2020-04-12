#ifndef EXPRESSIONS_HPP_
#define EXPRESSIONS_HPP_

#include "../json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json expressionCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2);

  if (val["ExpAct"].size() == 0) return val["falsey"][0];

  json _expVal = parser(val["ExpAct"], calc_params, vars, dir, false, line, true).exp["Hash_Values"];

  if (_expVal.find(index) == _expVal.end()) return _expVal["falsey"][0];
  else {
    json expVal = _expVal[index];

    indexes.erase(indexes.begin());

    if (indexes.size() == 0) return expVal[0];

    return indexesCalc(expVal[0]["Hash_Values"], indexes, calc_params, vars, line, dir, expVal[0]["Index_Type"]);
  }
}

#endif
