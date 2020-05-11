#ifndef MULTIPLYTYPES_HPP_
#define MULTIPLYTYPES_HPP_

#include <map>
#include <vector>
#include "../../values.hpp"
#include "../../json.hpp"
#include "../../structs.hpp"
using namespace std;
using json = nlohmann::json;

Action multiplystrings(Action num1, Action num2, json cli_params) {

  string fin = "";

  if (num1.Type == "number") {

    char* num = &num1.ExpStr[0][0];

    for (;((bool) IsLessC("0", num)); num = AddC(num, "1", &cli_params.dump()[0])) fin+=num2.ExpStr[0];

  } else {

    char* num = &num2.ExpStr[0][0];

    for (;((bool) IsLessC("0", num)); num = AddC(num, "1", &cli_params.dump()[0]))
      fin+=num1.ExpStr[0];
  }

  Action str = Action{ "string", "", { fin }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };

  char* i = "0";
  for (char it : fin) {

    Action character = strPlaceholder;

    character.ExpStr[0] = to_string(it);

    str.Hash_Values[string(i)] = { character };

    i = AddC(i, "1", &cli_params.dump()[0]);
  }

  return str;
}

Action multiplyarrays(Action num1, Action num2, json cli_params) {

  char* length = "0";

  //get length
  for (pair<string, vector<Action>> it : num1.Hash_Values) length = AddC(length, "1", &cli_params.dump()[0]);

  for (pair<string, vector<Action>> it : num2.Hash_Values) {

    string curIndex(AddC(length, &it.first[0], &cli_params.dump()[0]));

    num1.Hash_Values[curIndex] = it.second;
  }

  return num1;
}

#endif
