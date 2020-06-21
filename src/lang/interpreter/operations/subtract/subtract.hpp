#ifndef SUBTRACT_HPP_
#define SUBTRACT_HPP_

#include <map>
#include <vector>
#include "../../CliParams.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "subtracttypes.hpp"
#include "../numeric/modulo.hpp"
using CliParams = nlohmann::CliParams;

namespace omm {

  Action subtract(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num - num = num
      boolean - boolean = num
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num - num = num

      finalRet = subtractNums(num1, num2, cli_params);

    } else if (num1.Type == "boolean" && num2.Type == "bolean") { //detect case boolean - boolean = num

      finalRet = subtractbools(num1, num2, cli_params, this_vals, dir);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
