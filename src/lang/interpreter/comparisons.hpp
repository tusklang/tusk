#ifndef COMPARISONS_HPP_
#define COMPARISONS_HPP_

#include <map>
#include <deque>

#include "json.hpp"
#include "values.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

Action equals(Action val1, Action val2, json cli_params, map<string, Variable> vars, deque<map<string, vector<Action>>> this_vals, string dir) {

  if (val1.Name == "hashed_value" && val2.Name == "hashed_value") {

    //if the two hash values don't have the same size return false
    if (val1.Hash_Values.size() != val2.Hash_Values.size()) return falseRet;

    for (pair<string, vector<Action>> i : val1.Hash_Values) {

      auto finder = val2.Hash_Values.find(i.first);

      if (finder == val2.Hash_Values.end()) return falseRet;
      if (
        equals(
          parser(val2.Hash_Values[i.first], cli_params, vars, false, true, this_vals, dir).exp,
          parser(i.second, cli_params, vars, false, true, this_vals, dir).exp,
          cli_params,
          vars,
          this_vals,
          dir
        ).ExpStr[0] == "false"
        ) return falseRet;
    }

    return trueRet;
  } else {
    if (val1.ExpStr[0] == val2.ExpStr[0]) return trueRet;
    else return falseRet;
  }

  return falseRet;
}

Action isGreater(Action val1, Action val2, json cli_params) {
  if (val1.Type != "number" || val2.Type != "number") return falseRet;

  string num1 = val1.ExpStr[0]
  , num2 = val2.ExpStr[0];

  if (IsLessC(&num2[0], &num1[0])) return trueRet;

  return falseRet;
}

Action isLess(Action val1, Action val2, json cli_params) {
  if (val1.Type != "number" || val2.Type != "number") return falseRet;

  string num1 = val1.ExpStr[0]
  , num2 = val2.ExpStr[0];

  if (IsLessC(&num1[0], &num2[0])) return trueRet;

  return falseRet;

}

bool isTruthy(Action val) {
  return !(val.ExpStr[0] == "false" || val.Type == "falsey");
}

#endif
