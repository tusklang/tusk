#ifndef ADDTYPES_HPP_
#define ADDTYPES_HPP_

#include <map>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
using namespace std;
using json = nlohmann::json;

//file with all the functions to add different datatypes

Action addstrings(Action num1, Action num2, json cli_params) {

  string str = num1.ExpStr[0] + num2.ExpStr[0];

  map<string, vector<Action>> hash;

  char* i = "0";

  for (char c : str) {

    Action character = strPlaceholder;

    character.ExpStr[0] = to_string(c);
    hash[string(i)] = { character };

    i = AddC(i, "1", &cli_params.dump()[0]);
  }

  return Action{ "string", "", { str }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, hash, false };
}

Action addarrays(Action num1, Action num2, json cli_params) {

  map<string, vector<Action>> finalMap;

  if (num1.Type == "array") num1.Hash_Values[to_string(num1.Hash_Values.size())] = { num2 };
  else {

    for (pair<string, vector<Action>> it : num2.Hash_Values)
      finalMap[string(AddC(&it.first[0], "1", &cli_params.dump()[0]))] = { it.second };

    finalMap["0"] = { num1 };
  }

  return Action{ "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, finalMap, false };
}

Action addbools(Action num1, Action num2, json cli_params) {

  bool num1b = num1.ExpStr[0] == "true";
  bool num2b = num2.ExpStr[0] == "true";
  string final = to_string(num1b || num2b);

  return Action{ "boolean", "", { final }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };
}

Action addhashes(Action num1, Action num2, json cli_params) {

  map<string, vector<Action>> final = num1.Hash_Values;

  for (pair<string, vector<Action>> it : num2.Hash_Values) final[it.first] = it.second;

  return Action{ "hash", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 22, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, final, (num1.IsMutable || num2.IsMutable) };
}

#endif
