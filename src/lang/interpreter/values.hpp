#ifndef VALUES_HPP_
#define VALUES_HPP_

#include <vector>
#include <map>
#include <memory>
#include "structs.hpp"

namespace omm {

  std::vector<Action> emptyActVec;
  std::vector<std::vector<Action>> emptyActVec2D;
  std::map<std::string, std::vector<Action>> noneMap;
  std::vector<SubCaller> emptySubCaller;
  std::shared_ptr<std::future<Returner>> emptyFuture;

  //a bunch of commonly used values
  Action
    falseyVal = Action{ "falsey", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 41, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    trueRet = { "boolean", "", { "true" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    falseRet = { "boolean", "", { "false" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    zero = { "number", "", { "0" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    val1 = { "number", "", { "1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    valn1 = { "number", "", { "-1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    strPlaceholder = { "string", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    arrayVal = { "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture },
    hashVal = { "hash", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 22, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture };

}

#endif
