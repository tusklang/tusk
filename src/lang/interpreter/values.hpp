#ifndef VALUES_HPP_
#define VALUES_HPP_

#include <vector>
#include <map>
#include "structs.hpp"

namespace omm {

  std::vector<Action> emptyActVec;
  std::vector<std::vector<Action>> emptyActVec2D;
  std::map<std::string, std::vector<Action>> noneMap;

  //a bunch of commonly used values
  Action
    falseyVal = { "falsey", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 41, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    trueRet = { "boolean", "", { "true" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    falseRet = { "boolean", "", { "false" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    zero = { "number", "", { "0" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    val1 = { "number", "", { "1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    valn1 = { "number", "", { "-1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    strPlaceholder = { "string", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    arrayVal = { "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" },
    hashVal = { "hash", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 22, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };

}

#endif
