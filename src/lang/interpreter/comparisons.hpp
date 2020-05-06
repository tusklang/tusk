#ifndef COMPARISONS_HPP_
#define COMPARISONS_HPP_

#include "json.hpp"
#include "values.hpp"
using namespace std;

json equals(json val1, json val2, json cli_params, json vars) {

  if (val1["Name"] == "hashed_value" && val2["Name"] == "hashed_values") {

    if (val1["Hash_Values"].size() < val2["Hash_Values"].size()) {
      json temp = val1;

      val1 = val2;
      val2 = temp;
    }

    for (auto& i : val1["Hash_Values"].items()) {

      auto secondPart = val2["Hash_Values"].find(i.key());

      if (secondPart == val2["Hash_Values"].end()) return falseRet;
      if (parser(*secondPart, cli_params, vars, false, true).exp["ExpStr"][0] != parser(i.value(), cli_params, vars, false, true).exp["ExpStr"][0]) return falseRet;
    }

    return trueRet;
  } else {
    if (val1["ExpStr"][0] == val2["ExpStr"][0]) return trueRet;
    else return falseRet;
  }

  return falseRet;
}

json isGreater(json val1, json val2, json cli_params) {
  if (val1["Type"] != "number" || val2["Type"] != "number") return falseRet;

  string num1 = val1["ExpStr"][0]
  , num2 = val2["ExpStr"][0];

  if (IsLessC(&num2[0], &num1[0])) return trueRet;

  return falseRet;

}

json isLess(json val1, json val2, json cli_params) {
  if (val1["Type"] != "number" || val2["Type"] != "number") return falseRet;

  string num1 = val1["ExpStr"][0]
  , num2 = val2["ExpStr"][0];

  if (IsLessC(&num1[0], &num2[0])) return trueRet;

  return falseRet;

}

bool isTruthy(json val) {
  return !(val["ExpStr"][0] == "false" || val["Type"] == "falsey");
}

#endif
