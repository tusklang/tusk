#ifndef ADDTYPES_HPP_
#define ADDTYPES_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "../numeric/numeric.hpp"
using json = nlohmann::json;

namespace omm {

  //file with all the functions to add different datatypes

  Action addstrings(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::string str = num1.ExpStr[0] + num2.ExpStr[0];

    std::map<std::string, std::vector<Action>> hash;

    Action i = val1;

    for (char c : str) {

      Action character = strPlaceholder;

      character.ExpStr[0] = to_string(c);
      hash[normalize_number(i)] = { character };

      i = addNums(i, val1, cli_params);
    }

    return Action{ "string", "", { str }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, hash, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  }

  Action addarrays(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::map<std::string, std::vector<Action>> finalMap;

    if (num1.Type == "array") {
      num1.Hash_Values[to_string(num1.Hash_Values.size())] = { num2 };
      finalMap = num1.Hash_Values;
    } else {

      //re add when de-normalize is implemented
      // for (std::pair<std::string, std::vector<Action>> it : num2.Hash_Values)
      //   finalMap[string(addNums(&it.first[0], "1", &cli_params.dump()[0]))] = { it.second };

      finalMap["0"] = { num1 };
    }

    return Action{ "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, finalMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  }

  Action addbools(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    bool num1b = num1.ExpStr[0] == "true";
    bool num2b = num2.ExpStr[0] == "true";
    std::string final = to_string(num1b || num2b);

    return Action{ "boolean", "", { final }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  }

  Action addhashes(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::map<std::string, std::vector<Action>> final = num1.Hash_Values;

    for (std::pair<string, std::vector<Action>> it : num2.Hash_Values) final[it.first] = it.second;

    return Action{ "hash", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 22, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, final, (num1.IsMutable || num2.IsMutable), "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  }

}

#endif
