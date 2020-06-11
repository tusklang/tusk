#ifndef DIVIDE_HPP_
#define DIVIDE_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "../numeric/divide.hpp"
using json = nlohmann::json;

namespace omm {

  Action divide(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num / num = num
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num / num = num

      finalRet = divideNums(num1, num2, cli_params);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
