#ifndef OMM_NUMERIC_MODULO_HPP_
#define OMM_NUMERIC_MODULO_HPP_

#include <vector>

#include "../../structs.hpp"
#include "multiply.hpp"
#include "divide.hpp"
#include "normalize.hpp"

namespace omm {

  Action moduloNums(Action num1, Action num2, CliParams cli_params) {
    cli_params.Calc.PREC = 0; //force the number to be rounded down (cast to an int)
    Action divided = divideNums(num1, num2, cli_params);
    divided.Decimal.clear(); //round the number down
    Action multed = multiplyNums(divided, num2, cli_params); //multiply the divded value by the divisor
    return subtractNums(num1, multed, cli_params); //return quotient - product
  }

}

#endif
