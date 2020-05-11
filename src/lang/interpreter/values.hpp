#ifndef VALUES_HPP_
#define VALUES_HPP_

#include <vector>
#include <map>
#include "structs.hpp"
using namespace std;

vector<Action> emptyActVec;
vector<vector<Action>> emptyActVec2D;
map<string, vector<Action>> noneMap;

//a bunch of commonly used values
Action
  falseyVal = { "falsey", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 41, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  trueRet = { "boolean", "", { "true" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  falseRet = { "boolean", "", { "false" }, emptyActVec, {}, emptyActVec2D, {}, 40, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  zero = { "number", "", { "0" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  val1 = { "number", "", { "1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  valn1 = { "number", "", { "-1" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  strPlaceholder = { "string", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false },
  arrayVal = { "array", "", { "" }, emptyActVec, {}, emptyActVec2D, {}, 24, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };

#endif
