#ifndef OMM_NUMERIC_SUBTRACT_HPP_
#define OMM_NUMERIC_SUBTRACT_HPP_

#include "add.hpp"
#include "../../structs.hpp"
#include "../../json.hpp"
using json = nlohmann::json;

namespace omm {
  Action subtract(Action num1, Action num2, json cli_params) {

    //negate the num2
    for (int i = 0; i < num2.Decimal.size(); ++i) num2.Decimal[i]*=-1;
    for (int i = 0; i < num2.Integer.size(); ++i) num2.Integer[i]*=-1;

    return add(num1, num2, cli_params);
  }
}

#endif
