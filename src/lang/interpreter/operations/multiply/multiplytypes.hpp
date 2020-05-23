#ifndef MULTIPLYTYPES_HPP_
#define MULTIPLYTYPES_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../values.hpp"
#include "../../json.hpp"
#include "../../structs.hpp"
using json = nlohmann::json;

namespace omm {

  Action multiplystrings(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::string fin = "";

    if (num1.Type == "number") {

      char* num = &num1.ExpStr[0][0];

      for (;((bool) IsLessC("0", num)); num = AddC(num, "1", &cli_params.dump()[0])) fin+=num2.ExpStr[0];

    } else {

      char* num = &num2.ExpStr[0][0];

      for (;((bool) IsLessC("0", num)); num = AddC(num, "1", &cli_params.dump()[0]))
        fin+=num1.ExpStr[0];
    }

    Action str = Action{ "string", "", { fin }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };

    char* i = "0";
    for (char it : fin) {

      Action character = strPlaceholder;

      character.ExpStr[0] = to_string(it);

      str.Hash_Values[string(i)] = { character };

      i = AddC(i, "1", &cli_params.dump()[0]);
    }

    return str;
  }

  Action multiplyarrays(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    char* length = "0";

    //get length
    for (std::pair<std::string, std::vector<Action>> it : num1.Hash_Values) {

      //skip over the falsey index
      if (it.first == "falsey") continue;

      length = AddC(length, "1", &cli_params.dump()[0]);
    }

    for (std::pair<std::string, std::vector<Action>> it : num2.Hash_Values) {

      if (it.first == "falsey") continue;

      std::string curIndex(AddC(length, &it.first[0], &cli_params.dump()[0]));

      num1.Hash_Values[curIndex] = it.second;
    }

    return num1;
  }

}

#endif
