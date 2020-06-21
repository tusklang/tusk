#ifndef OMM_NUMERIC_ADD_HPP_
#define OMM_NUMERIC_ADD_HPP_

#include <vector>
#include <string>

#include "utils.hpp"
#include "../../values.hpp"
#include "../../structs.hpp"

namespace omm {

  Action addNums(Action num1, Action num2, CliParams cli_params) {

    /*

    Example:

        1  2  3  4
     +    -5 -7 -9
      -------------
        1 -3 -4 -5

      THOSE ARE THE DIGITS
    */

    std::vector<long long>
      int1 = num1.Integer,
      int2 = num2.Integer,
      dec1 = num1.Decimal,
      dec2 = num2.Decimal;

    //ensure that int1 and dec1 have greater lengths
    if (int1.size() < int2.size()) {

      //swap the values
      std::vector<long long> temp = int1;
      int1 = int2;
      int2 = temp;
    }
    if (dec1.size() < dec2.size()) {

      //swap the values
      std::vector<long long> temp = dec1;
      dec1 = dec2;
      dec2 = temp;
    }

    int carry = 0;

    //got this from the python source code
    //https://github.com/python/cpython/blob/master/Objects/longobject.c#L3003
    //basically, it does not require an "if" statement to see if num2 has an element at an index

    std::vector<long long> newDecimal(dec1.size());

    //do the decimal
    int decIndex;

    for (decIndex = 0; decIndex < dec2.size(); ++decIndex) {

      long long summed = dec2[decIndex] + dec1[decIndex] + carry;
      carry = 0;

      if (summed > OMM_MAX_DIGIT) {
        summed-=OMM_MAX_DIGIT;
        carry = 1;
      }
      if (summed < OMM_MIN_DIGIT) {
        summed+=OMM_MIN_DIGIT;
        carry = -1;
      }

      newDecimal[decIndex] = summed;
    }
    for (;decIndex < dec1.size(); ++decIndex) {

      long long summed = dec1[decIndex] + carry;
      carry = 0;

      if (summed > OMM_MAX_DIGIT) {
        summed-=OMM_MAX_DIGIT;
        carry = 1;
      }
      if (summed < OMM_MIN_DIGIT) {
        summed+=OMM_MIN_DIGIT;
        carry = -1;
      }

      newDecimal[decIndex] = summed;
    }

    std::vector<long long> newInteger(int1.size());

    //do the integer
    int intIndex;

    for (intIndex = 0; intIndex < int2.size(); ++intIndex) {

      long long summed = int2[intIndex] + int1[intIndex] + carry;
      carry = 0;

      if (summed >= OMM_MAX_DIGIT) {
        summed-=OMM_MAX_DIGIT;
        carry = 1;
      }
      if (summed <= OMM_MIN_DIGIT) {
        summed+=OMM_MIN_DIGIT;
        carry = -1;
      }

      newInteger[intIndex] = summed;
    }
    for (;intIndex < int1.size(); ++intIndex) {

      long long summed = int1[intIndex] + carry;
      carry = 0;

      if (summed >= OMM_MAX_DIGIT) {
        summed-=OMM_MAX_DIGIT;
        carry = 1;
      }
      if (summed <= OMM_MIN_DIGIT) {
        summed+=OMM_MIN_DIGIT;
        carry = -1;
      }

      newInteger[intIndex] = summed;
    }

    newInteger.push_back(carry);

    while (newInteger[newInteger.size() - 1] == 0) newInteger.pop_back(); //remove leading zeros

    return Action{ "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, newInteger, newDecimal, emptyFuture };
  }

}

#endif
