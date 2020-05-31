#ifndef MULTIPLY_HPP_
#define MULTIPLY_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "multiplytypes.hpp"
using json = nlohmann::json;

namespace omm {

  Action multiply(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num * num = num
      string * num = string
      array * array = array
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num * num = num

      std::string val(MultiplyC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]));

      finalRet = Action{ "number", "", { val }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture };

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
