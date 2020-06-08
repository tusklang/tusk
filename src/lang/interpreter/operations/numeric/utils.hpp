#ifndef OMM_NUMERIC_OPERATIONS_UTILS_HPP_
#define OMM_NUMERIC_OPERATIONS_UTILS_HPP_

#include <iostream>
#include <vector>
#include <algorithm>
#include <cmath>

#include "../../structs.hpp"
#include "../../json.hpp"
#include "../../values.hpp"
#include "subtract.hpp"
#include "multiply.hpp"
#include "normalize.hpp"
using json = nlohmann::json;

namespace omm {

  bool isLess(Action num1, Action num2, json cli_params) {

    bool swappedInt = false /* if a swap was performed (for integer) */, swappedDec = false /* if a swap was performed (for decimal) */;

    if (num1.Integer.size() < num2.Integer.size()) { //determine swappedInt
      swappedInt = true;

      //swap num1 and num2
      std::vector<long long> temp = num2.Integer;
      num2.Integer = num1.Integer;
      num1.Integer = temp;
    }
    if (num1.Decimal.size() < num2.Decimal.size()) { //determine swappedDec
      swappedDec = true;

      //swap num1 and num2
      std::vector<long long> temp = num2.Decimal;
      num2.Decimal = num1.Decimal;
      num1.Decimal = temp;
    }

    //do the integer
    int intI;

    //do for num1
    for (intI = num1.Integer.size() - 1; intI >= num2.Integer.size(); --intI)
      if (num1.Integer[intI] < 0) return !swappedInt;
      else if (num1.Integer[intI] > 0) return swappedInt;
    //do for num2
    for (;intI >= 0; --intI)
      if (num1.Integer[intI] < num2.Integer[intI]) return !swappedInt;
      else if (num1.Integer[intI] > num2.Integer[intI]) return swappedInt;

    if (num1.Decimal.size() == 0) return false; //safeguard to make sure there is a decimal

    //do the decimal
    int decI;

    //do for num2
    for (decI = num1.Decimal.size() - 1; decI >= num2.Decimal.size(); --decI)
      if (num1.Decimal[decI] < 0) return !swappedDec;
      else if (num1.Decimal[decI] > 0) return swappedDec;
    //do for num2
    for (;decI >= 0; --decI)
      if (num1.Decimal[decI] < num2.Integer[decI]) return !swappedDec;
      else if (num1.Decimal[decI] > num2.Integer[decI]) return swappedDec;


    return false;
  }

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
