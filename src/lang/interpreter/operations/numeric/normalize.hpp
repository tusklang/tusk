#ifndef OMM_NUMBER_NORMALIZE_HPP_
#define OMM_NUMBER_NORMALIZE_HPP_

#include <cstdlib>
#include <vector>

#include "../../structs.hpp"
#include "utils.hpp"

namespace omm {

  //function to convert integer.decimal to a string
  std::string normalize_number(Action num) {

    /*

    ALGORITHM TO NORMALIZE:

        Starting with this number:
          [3412, -9912, 0001]
          (in omm this is stored as [0001, -9912, 3412])

        STEP 0: (initializer step)
          remove the leading zeros in the integer

          if the first digit non zero digit is negative set isNeg = true, otherwise isNeg = false
          (if there is no first digit, go to the decimals and see)

        STEP 1:
          loop through each number (from from decimal to integer)
          in each iteration of the loop, if the number is the opposite of `isNeg` (meaning if `isNeg` is false, then the current value should be positive and vice versa)
          use the following expression to get the complement

            `OMM_MAX_DIGIT` - |`current num`|

          replace the `current num` with this new value.
          next, if `isNeg`, the digit to the left should be added by one, otherwise, subtract it by 1.
          go to the next value and repeat

        STEP 2:
          join the vector of the integer and decimal with '.', then join each digit with ''
          if `isNeg` then precede the string with a '-'
          finally, return the result

    */

    std::vector<long long>
      integer = num.Integer,
      decimal = num.Decimal;

    //the first digit is actually the last index
    //because omm numbers are stored as so [1234, 5678, 9101] = 910, 156, 781, 234

    //remove leading zeros in the integer
    while (integer.size() != 0 && integer[integer.size() - 1] == 0) integer.pop_back();

    //detect if it is negative
    bool isNeg = false;

    if (integer.size() == 0) {
      if (decimal.size() == 0) return "0"; //this means that there is no number

      int decIndexCounter;

      for (decIndexCounter = decimal.size(); decimal.size() - 1 != -1 && decimal[decimal.size() - 1] == 0; decimal.pop_back());

      if (decimal.size() - 1 == -1) return "0";

      if (decimal[decIndexCounter - 1] < 0) isNeg = true;
    } else if (integer[integer.size() - 1] < 0) isNeg = true;

    int carry = 0;

    for (int i = 0; i < decimal.size(); ++i) {
      bool curIsNeg = decimal[i] < 0; //get if it is negative
      decimal[i]+=carry;

      decimal[i] = std::abs(decimal[i]);

      if (curIsNeg != isNeg && decimal[i] != 0 /* prevent zeros from being counted */ ) {
        decimal[i] = OMM_MAX_DIGIT - decimal[i];
        carry = isNeg ? 1 : -1;
        continue;
      }

      carry = 0;

    }

    for (int i = 0; i < integer.size(); ++i) {
      bool curIsNeg = integer[i] < 0; //get if it is negative
      integer[i]+=carry;

      integer[i] = std::abs(integer[i]);

      if (curIsNeg != isNeg && integer[i] != 0 /* prevent zeros from being counted */ ) {
        integer[i] = OMM_MAX_DIGIT - integer[i];
        carry = isNeg ? 1 : -1;
        continue;
      }

      carry = 0;

    }

    //combine for final result
    std::string joined = "";

    if (decimal.size() != 0) { //this is because if there is no decimal, a "." will still be inserted
      for (long long it : decimal)
        joined = std::to_string(it) + joined;
      joined = "." + joined;
    }

    for (long long it : integer)
      joined = std::to_string(it) + joined;

    joined = (isNeg ? "-" : "") + joined;

    return joined;
  }

}

#endif
