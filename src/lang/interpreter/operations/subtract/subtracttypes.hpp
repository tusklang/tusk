#ifndef SUBTRACTTYPES_HPP_
#define SUBTRACTTYPES_HPP_

#include <map>
#include <vector>
#include "../../values.hpp"
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../comparisons.hpp"
#include "../../parser.hpp"
using namespace std;
using json = nlohmann::json;

Returner parser(const vector<Action> actions, const json cli_params, map<string, Variable> vars, const bool groupReturn, const bool expReturn);

Action subtractstrings(Action num1, Action num2, json cli_params) {

  map<string, vector<Action>> finalMap;

  if (num1.Type == "number") {

    for (pair<string, vector<Action>> it : num2.Hash_Values) {
      string k = it.first;
      vector<Action> v = it.second;

      if (((bool) IsLessC(&k[0], &num1.ExpStr[0][0])) || string(ReturnInitC(&num1.ExpStr[0][0])) == string(ReturnInitC(&k[0]))) continue;

      string cur(SubtractC(&num1.ExpStr[0][0], &k[0], &cli_params.dump()[0]));

      finalMap[cur] = v;
    }

  } else {

    char* length = "0";

    //get string length
    for (pair<string, vector<Action>> it : num1.Hash_Values)
      length = AddC(length, "1", &cli_params.dump()[0]);

    while ((bool) IsLessC("0", &num2.ExpStr[0][0])) {

      string cur(AddC(length, &num2.ExpStr[0][0], &cli_params.dump()[0]));
      num1.Hash_Values.erase(cur);

      string subtracted(SubtractC(&num2.ExpStr[0][0], "1", &cli_params.dump()[0]));
      num2.ExpStr[0] = subtracted;
    }

    finalMap = num1.Hash_Values;
  }

  string str;
  map<string, Variable> emptyVars;

  for (pair<string, vector<Action>> it : finalMap)
    str+=parser(it.second, cli_params, emptyVars, false, true).exp.ExpStr[0];

  return Action{ "string", "", { str }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, finalMap, false };
}

Action subtractbools(Action num1, Action num2, json cli_params) {

  bool num1Bool = num1.ExpStr[0] == "true";
  bool num2Bool = num2.ExpStr[0] == "true";
  string calc = to_string(((int) num1Bool) - ((int) num2Bool));

  return Action{ "boolean", "", { calc }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };
}

Action subtractarrays(Action num1, Action num2, json cli_params) {

  map<string, vector<Action>> hash = subtractstrings(num1, num2, cli_params).Hash_Values;

  return Action{ "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, hash, false };
}

#endif
