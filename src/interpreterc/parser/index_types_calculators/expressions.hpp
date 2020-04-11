#ifndef EXPRESSIONS_HPP_
#define EXPRESSIONS_HPP_

#include "../json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json expressionCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2)
  , strVal = parser(val, calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

  strVal = strVal.substr(1, strVal.length() - 2);

  //OPTIMIZE:

  while (IsLessC("0", &index[0])) {
    strVal = strVal.substr(1);
    index = Subtract(&index[0], "1", &calc_params.dump()[0], line);
  }

  json type = {
    {"Type", "string"},
    {"Name", ""},
    {"ExpStr", json::parse("[\"\'" + string(1, strVal[0]) + "\'\"]")},
    {"ExpAct", "[]"_json},
    {"Params", "[]"_json},
    {"Args", "[]"_json},
    {"Condition", "[]"_json},
    {"ID", 38},
    {"First", "[]"_json},
    {"Second", "[]"_json},
    {"Degree", "[]"_json},
    {"Value", "[[]]"_json},
    {"Indexes", "[[]]"_json},
    {"Index_Type", ""},
    {"Hash_Values", "{}"_json},
    {"ValueType", "[]"_json}
  }
  , returner = {
    {"Type", "string"},
    {"Name", ""},
    {"ExpStr", json::parse("[\"\'" + string(1, strVal[0]) + "\'\"]")},
    {"ExpAct", "[]"_json},
    {"Params", "[]"_json},
    {"Args", "[]"_json},
    {"Condition", "[]"_json},
    {"ID", 38},
    {"First", "[]"_json},
    {"Second", "[]"_json},
    {"Degree", "[]"_json},
    {"Value", "[[]]"_json},
    {"Indexes", "[[]]"_json},
    {"Index_Type", ""},
    {"Hash_Values", "{}"_json},
    {"ValueType", json::parse("[" + type.dump() + "]")}
  };

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return returner;

  return indexesCalc(returner["Index_Type"] == "hash" ? returner["Hash_Values"] : (returner["Index_Type"] == "array" ? returner["Value"] : returner["ExpAct"]), indexes, calc_params, vars, line, dir, returner["Index_Type"]);
}

#endif
