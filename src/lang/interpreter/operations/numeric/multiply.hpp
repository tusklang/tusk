#ifndef OMM_NUMERIC_MULTIPLY_HPP_
#define OMM_NUMERIC_MULTIPLY_HPP_

#include <vector>

#include "utils.hpp"
#include "add.hpp"
#include "../../json.hpp"
#include "../../values.hpp"
using json = nlohmann::json;

namespace omm {

  //using naive approach because karatsuba makes digits >= 10
  //example
  /*
    24
  x 19
  -----

  2 * 1 = 2
  4 * 9 = 36 <-- is >= 10

  omm has numbers that go from -9999 to 9999 and karatsuba allows for numbers >= 10000
  */
  Action multiply(Action num1, Action num2, json cli_params) {

    std::vector<std::vector<long long>> multFin; //store the final values that were multiplied
    int trailingZeroCount = 0;

    //amount of decimal places there are
    int decPlaceCount = num1.Decimal.size() + num2.Decimal.size();

    std::vector<long long> num1n = num1.Decimal, num2n = num2.Decimal; //actual numbers for num1 and num2

    num1n.insert(num1n.end(), num1.Integer.begin(), num1.Integer.end());
    num2n.insert(num2n.end(), num2.Integer.begin(), num2.Integer.end());

    if (num1n.size() < num2n.size()) { //swap num1n and num2n if num2n is greater (improves performance)
      std::vector<long long> temp = num1n;
      num1n = num2n;
      num2n = temp;
    }

    for (long long i : num1n) {

      multFin.push_back({}); //add a new row to the final multiplied values

      long long carry = 0; //variable to carry over overflowed numbers

      for (int i = 0; i < trailingZeroCount; ++i) multFin[multFin.size() - 1].push_back(0); //insert the trailing zeros

      for (long long o : num2n) {

        long long product = i * o + carry;

        carry = 0; //reverrt carry after it was factored in

        //determine if a carry is neccessary
        if (product >= OMM_MAX_DIGIT) {

          long long rounded = product / OMM_MAX_DIGIT * OMM_MAX_DIGIT; //round down by `OMM_MAX_DIGIT`
          carry = product / OMM_MAX_DIGIT; //divide by the max digit to get the carry
          product-=rounded;

        }
        if (product <= OMM_MIN_DIGIT) {
          long long rounded = ((product + OMM_MIN_DIGIT - 1) / OMM_MIN_DIGIT) * OMM_MIN_DIGIT; //round up by `OMM_MIN_DIGIT`
          carry = product / OMM_MAX_DIGIT; //divide by the max digit to get the carry (not min digit)
          product-=rounded;
        }

        multFin[multFin.size() - 1].push_back(product);
      }

      multFin[multFin.size() - 1].push_back(carry); //add the carry

      trailingZeroCount++;
    }

    std::vector<long long> totalSum = { 0 }; //the summed value

    for (std::vector<long long> i : multFin) {

      Action totalSumAct = Action{ "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, totalSum, emptyLLVec, emptyFuture }; //placeholder number to pass into add
      Action multFinAct = Action{ "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, i, emptyLLVec, emptyFuture }; //placeholder number to pass into add

      totalSum = add(totalSumAct, multFinAct, cli_params).Integer;
    }

    //get the decimal and the integer
    std::vector<long long>
      decimalRet(totalSum.begin(), totalSum.begin() + decPlaceCount),
      integerRet(totalSum.begin() + decPlaceCount, totalSum.end());

    return Action{ "number", "", {}, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, integerRet, decimalRet, emptyFuture };
  }

}

#endif
