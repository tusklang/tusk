#ifndef MULTIPLYTYPES_HPP_
#define MULTIPLYTYPES_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../values.hpp"
#include "../../CliParams.hpp"
#include "../../structs.hpp"
#include "../numeric/numeric.hpp"
using CliParams = nlohmann::CliParams;

namespace omm {

  Action multiplystrings(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    std::string fin = "";

    if (num1.Type == "number") {

      Action num = num1;

      for (;(isLess(zero, num, cli_params)); num = addNums(num, val1, cli_params)) fin+=num2.ExpStr[0];

    } else {

      Action num = num2;

      for (;(isLess(zero, num, cli_params)); num = addNums(num, val1, cli_params))
        fin+=num1.ExpStr[0];
    }

    Action str = Action{ "string", "", { fin }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };

    Action i = zero;
    for (char it : fin) {

      Action character = strPlaceholder;

      character.ExpStr[0] = to_string(it);

      str.Hash_Values[normalize_number(i)] = { character };

      i = addNums(i, val1, cli_params);
    }

    return str;
  }

  Action multiplyarrays(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    Action length = zero;

    //get length
    for (std::pair<std::string, std::vector<Action>> it : num1.Hash_Values) {

      //skip over the falsey index
      if (it.first == "falsey") continue;

      length = addNums(length, val1, cli_params);
    }

    for (std::pair<std::string, std::vector<Action>> it : num2.Hash_Values) {

      if (it.first == "falsey") continue;

      //fix later (when swig is being used) (using val1 for now)
      Action curIndex = addNums(length, val1 /*it.first*/, cli_params);

      num1.Hash_Values[normalize_number(curIndex)] = it.second;
    }

    return num1;
  }

}

#endif
