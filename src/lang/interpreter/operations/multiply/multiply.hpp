#ifndef MULTIPLY_HPP_
#define MULTIPLY_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../structs.hpp"
#include "../../values.hpp"
#include "multiplytypes.hpp"
#include "../numeric/modulo.hpp"

namespace omm {

  Action multiply(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num * num = num
      string * num = string
      array * array = array
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num * num = num

      finalRet = multiplyNums(num1, num2, cli_params);

    } else if ((num1.Type == "string" && num2.Type == "number") || (num1.Type == "number" && num2.Type == "string")) { //detect case string * num = string

      finalRet = multiplystrings(num1, num2, cli_params, this_vals, dir);

    } else if (num1.Type == "array" && num2.Type == "array") { //detect case array * array = array

      finalRet = multiplyarrays(num1, num2, cli_params, this_vals, dir);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
