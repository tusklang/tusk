#ifndef OMM_NUMERIC_OPERATIONS_UTILS_HPP_
#define OMM_NUMERIC_OPERATIONS_UTILS_HPP_

#include <vector>
#include <algorithm>
#include <cmath>

#include "../../structs.hpp"

namespace omm {

  //the max digit number (if digit surpasses this, overflow it)

  const int DigitSize = 4; //if you change here, change in numconv.go

  const long long
    OMM_MAX_DIGIT = std::pow(10, DigitSize), //the actual max digit + 1
    OMM_MIN_DIGIT = -1 * OMM_MAX_DIGIT; //the actual min digit - 1

  namespace numeric_utils {

    bool isLess(Action num1Act, Action num2Act) {

      std::vector<long long> num1 = num1Act.Decimal, num2 = num2Act.Decimal;

      //merge the decimals with the integers
      num1.insert(num1.end(), num1Act.Integer.begin(), num1Act.Integer.end());
      num2.insert(num2.end(), num2Act.Integer.begin(), num2Act.Integer.end());

      //reverse the vectors
      std::reverse(num1.begin(), num1.end());
      std::reverse(num2.begin(), num2.end());

      for (int i = 0; i < num1.size() && i < num2.size(); ++i) {

        long long val1, val2;

        //account for missing values
        if (num1.size() <= i) val1 = 0;
        else val1 = num1[i];

        if (num2.size() <= i) val2 = 0;
        else val2 = num2[i];

        if (val1 < val2) return true;
        else if (val1 > val2) break;

      }

      return false;
    }

  }

}

#endif
