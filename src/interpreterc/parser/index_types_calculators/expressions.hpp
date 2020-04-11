#ifndef EXPRESSIONS_HPP_
#define EXPRESSIONS_HPP_

#include "../json.hpp"
using namespace std;

json indexesCalc(json val, json indexes, json calc_params, json vars, int line, string dir, string type);

json expressionCalc(json val, json indexes, json calc_params, json vars, int line, string dir) {

  json _index = parser(indexes[0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
  string index = _index.dump().substr(1, _index.dump().length() - 2);

  string strVal = parser(val, calc_params, vars, line, false, dir, true).exp["ExpStr"][0].get<string>();

  //OPTIMIZE:

  while (IsLessC("0", &index[0])) {
    strVal = strVal.substr(1);
    index = Subtract(&index[0], "1", &calc_params.dump()[0], line);
  }

  if (indexes.size() == 0) return {
    {"Type", "string"},
    {"Name", ""},
    {"ExpStr", json::parse("[" + strVal + "]")},
    {"ExpAct", "[]"_json},
    {"Params", "[]"_json},
    {"Args", "[]"_json},
    {"Condition", "[]"_json},
    {"ID", 38},
    {"First", "[]"_json},
    {"Second", "[]"_json},
    {"Degree", "[]"_json},
    {"Value", "[[]]"_json},
    {"Indexes", "[[]]"},
    {"Index_Type", ""},
    {"Hash_Values", "{}"_json},
    {"ValueType", {{
        {"Type", "string"},
        {"Name", ""},
        {"ExpStr", json::parse("[" + strVal + "]")},
        {"ExpAct", "[]"_json},
        {"Params", "[]"_json},
        {"Args", "[]"_json},
        {"Condition", "[]"_json},
        {"ID", 38},
        {"First", "[]"_json},
        {"Second", "[]"_json},
        {"Degree", "[]"_json},
        {"Value", "[[]]"_json},
        {"Indexes", "[[]]"},
        {"Index_Type", ""},
        {"Hash_Values", "{}"_json},
        {"ValueType", "[]"_json}
      }}
    }
  };

  return indexesCalc(indexVal[0]["Index_Type"] == "hash" ? indexVal[0]["Hash_Values"] : (indexVal[0]["Index_Type"] == "array" ? indexVal[0]["Value"] : indexVal[0]), indexes, calc_params, vars, line, dir, indexVal[0]["Index_Type"]);
}

#endif
