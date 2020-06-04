#ifndef VALUES_HPP_
#define VALUES_HPP_

#include <vector>
#include <map>
#include <memory>
#include "structs.hpp"

namespace omm {

  #define A (DigitSize);

  std::vector<Action> emptyActVec;
  std::vector<long long> emptyLLVec;
  std::vector<std::vector<Action>> emptyActVec2D;
  std::map<std::string, std::vector<Action>> noneMap;
  std::vector<SubCaller> emptySubCaller;
  std::shared_ptr<std::future<Returner>> emptyFuture;

  //a bunch of commonly used values
  Action
    falseyVal = Action{ "falsey", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 41, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture },
    trueRet = { "boolean", "", { "true" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture },
    falseRet = { "boolean", "", { "false" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture },
    zero = { "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, { 0 }, emptyLLVec, emptyFuture },
    val1 = { "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, { 1001, 1001 }, {}, emptyFuture },
    valn1 = { "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, { -1 }, emptyLLVec, emptyFuture },
    strPlaceholder = { "string", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture },
    arrayVal = { "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture },
    hashVal = { "hash", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 22, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };

}

#endif
