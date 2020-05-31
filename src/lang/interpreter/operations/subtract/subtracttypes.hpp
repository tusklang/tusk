#ifndef SUBTRACTTYPES_HPP_
#define SUBTRACTTYPES_HPP_

#include <map>
#include <vector>
#include <deque>
#include "../../values.hpp"
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../comparisons.hpp"
#include "../../parser.hpp"
using json = nlohmann::json;

namespace omm {

  Returner parser(const std::vector<Action> actions, const json cli_params, std::map<std::string, Variable> vars, const bool groupReturn, const bool expReturn, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir);

  Action subtractstrings(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::map<string, std::vector<Action>> finalMap;

    if (num1.Type == "number") {

      for (std::pair<std::string, std::vector<Action>> it : num2.Hash_Values) {
        std::string k = it.first;
        std::vector<Action> v = it.second;

        if (((bool) IsLessC(&k[0], &num1.ExpStr[0][0])) || string(ReturnInitC(&num1.ExpStr[0][0])) == string(ReturnInitC(&k[0]))) continue;

        std::string cur(SubtractC(&num1.ExpStr[0][0], &k[0], &cli_params.dump()[0]));

        finalMap[cur] = v;
      }

    } else {

      char* length = "0";

      //get string length
      for (std::pair<std::string, std::vector<Action>> it : num1.Hash_Values)
        length = AddC(length, "1", &cli_params.dump()[0]);

      while ((bool) IsLessC("0", &num2.ExpStr[0][0])) {

        std::string cur(AddC(length, &num2.ExpStr[0][0], &cli_params.dump()[0]));
        num1.Hash_Values.erase(cur);

        std::string subtracted(SubtractC(&num2.ExpStr[0][0], "1", &cli_params.dump()[0]));
        num2.ExpStr[0] = subtracted;
      }

      finalMap = num1.Hash_Values;
    }

    std::string str;
    std::map<std::string, Variable> emptyVars;

    for (std::pair<std::string, std::vector<Action>> it : finalMap)
      str+=parser(it.second, cli_params, emptyVars, false, true, this_vals, dir).exp.ExpStr[0];

    return Action{ "string", "", { str }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, finalMap, false, "private", emptySubCaller, emptyFuture };
  }

  Action subtractbools(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    bool num1Bool = num1.ExpStr[0] == "true";
    bool num2Bool = num2.ExpStr[0] == "true";
    std::string calc = to_string(((int) num1Bool) - ((int) num2Bool));

    return Action{ "boolean", "", { calc }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture };
  }

  Action subtractarrays(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::map<std::string, std::vector<Action>> hash = subtractstrings(num1, num2, cli_params, this_vals, dir).Hash_Values;

    return Action{ "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, hash, false, "private", emptySubCaller, emptyFuture };
  }

}

#endif
