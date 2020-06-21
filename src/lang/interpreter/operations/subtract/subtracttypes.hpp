#ifndef SUBTRACTTYPES_HPP_
#define SUBTRACTTYPES_HPP_

#include <map>
#include <vector>
#include <deque>
#include "../../values.hpp"
#include "../../CliParams.hpp"
#include "../../structs.hpp"
using CliParams = nlohmann::CliParams;

namespace omm {

  Action subtractbools(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    bool num1Bool = num1.ExpStr[0] == "true";
    bool num2Bool = num2.ExpStr[0] == "true";
    std::string calc = to_string(((int) num1Bool) - ((int) num2Bool));

    return Action{ "boolean", "", { calc }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  }

}

#endif
