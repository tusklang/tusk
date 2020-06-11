#ifndef OMM_NUMERIC_DIVIDE_HPP_
#define OMM_NUMERIC_DIVIDE_HPP_

#include <vector>
#include <deque>

#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "utils.hpp"
#include "subtract.hpp"
#include "multiply.hpp"
using json = nlohmann::json;

namespace omm {

  //using long division algorithm
  Action divideNums(Action num1, Action num2, json cli_params) {

    //maybe in a future version switch to the algorithm python uses
    //https://github.com/python/cpython/blob/8bd216dfede9cb2d5bedb67f20a30c99844dbfb8/Objects/longobject.c#L2610
    //because it is faster

    //num2 is the divisor
    //num1 is the dividend

    if (equals(num2, zero, cli_params)) return falseyVal; //if it is n/0 return undef
    if (equals(num1, zero, cli_params)) return zero; //if it is 0/n return zero

    int addDecimalPlaces = num1.Integer.size() + num2.Decimal.size(); //the amount of decimal places in the final result

    std::vector<long long> num1n = num1.Integer;

    //merge the integer with the decimal
    num1n.insert(num1n.end(), num1.Decimal.begin(), num1.Decimal.end());

    std::deque<long long> num1D(num1n.begin(), num1n.end()); //convert the vector to a deque (to push to front)

    //manage precision
    while (num2.Integer.size() + num2.Decimal.size() > num1D.size()) num1D.push_front(0);

    int prec = cli_params["Calc"]["PREC"].get<int>();

    while (num1D.size() < prec) num1D.push_front(0);

    std::vector<long long> temp(num1D.begin(), num1D.end()); //copy thr deque to a vector
    num1n = temp;

    std::deque<long long> curVal; //current value under the "house" of the long division
    std::deque<long long> final;

    Action num2Abs = abs(num2, cli_params);

    //do the actual division
    for (std::vector<long long>::reverse_iterator it = num1n.rbegin(); it != num1n.rend(); ++it) {
      long long i = *it; //get the value

      curVal.push_front(i);

      Action curValAct = zero;
      curValAct.Integer = std::vector<long long>(curVal.begin(), curVal.end());

      Action curValAbs = abs(curValAct, cli_params);

      //if the |current value| is less than the |num2|
      if (isLess(curValAbs, num2Abs, cli_params)) {
        final.push_front(0);
        continue;
      }

      Action
        curQuotient = zero, //the quotient of the current iteration
        added = zero;

      for (Action addedTemp = added; isLess(addedTemp = addNums(added, num2Abs, cli_params), curValAbs, cli_params) || equals(addedTemp, curValAbs, cli_params); added = addedTemp)
        curQuotient = addNums(curQuotient, val1, cli_params); //increment the current quotient

      if (isLess(curValAct, zero, cli_params)) //detect a negative current quotient
        curQuotient = multiplyNums(curQuotient, valn1, cli_params);

      std::vector<long long> subtractedVec = subtractNums(curValAbs, added, cli_params).Integer;
      curVal = std::deque<long long>(subtractedVec.begin(), subtractedVec.end());

      final.insert(final.begin(), curQuotient.Integer.begin(), curQuotient.Integer.end());
    }

    //seperate the decimal with the integer
    std::vector<long long>
      decimal(final.begin(), final.end() - addDecimalPlaces),
      integer(final.end() - addDecimalPlaces, final.end());

    Action finalAct = zero;
    finalAct.Integer = integer;
    finalAct.Decimal = decimal;

    return finalAct;
  }

}

#endif
