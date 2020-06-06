#ifndef OMM_NUMERIC_OPERATIONS_UTILS_HPP_
#define OMM_NUMERIC_OPERATIONS_UTILS_HPP_

#include <vector>
#include <algorithm>
#include <cmath>

#include "../../structs.hpp"
#include "../../json.hpp"
#include "../../values.hpp"
#include "subtract.hpp"
#include "multiply.hpp"
using json = nlohmann::json;

namespace omm {

  bool isLessVec(std::vector<long long> num1, std::vector<long long> num2, bool isDec) { //helper func to isLess to determine if one vector is less than another

    std::vector<long long> greater = num1.size() > num2.size() ? num1 : num2; //the number with the greater length

    for (int i = greater.size(); i >= 0; --i) {
      long long i1, i2; //declare the vars

      //account for missing values
      if (i >= num1.size()) i1 = 0;
      else i1 = num1[i];
      if (i >= num2.size()) i2 = 0;
      else i2 = num2[i];

      if (i1 < i2) return true;
      else if (i1 > i2) return false;

    }

    return false;
  }

  bool isLess(Action num1, Action num2, json cli_params) {

    bool
      intLess = isLessVec(num1.Integer, num2.Integer, false),
      decLess = isLessVec(num1.Decimal, num2.Decimal, true);

    return intLess ? true : decLess;
  }

  Action subtractNums(Action num1, Action num2, json cli_params);

  bool equals(Action num1Act, Action num2Act, json cli_params) {

    //uses num1 - num2 == 0

    Action subtracted = subtractNums(num1Act, num2Act, cli_params);

    for (long long i : subtracted.Integer)
      if (i != 0) return false;
    for (long long i : subtracted.Decimal)
      if (i != 0) return false;

    return true;
  }

  bool isTruthy(Action val) {
    return !(val.ExpStr[0] == "false" || val.Type == "falsey");
  }

  Action abs(Action val, json cli_params) {
    if (isLess(val, zero, cli_params)) return multiplyNums(val, valn1, cli_params);

    return val;
  }

}

#endif
