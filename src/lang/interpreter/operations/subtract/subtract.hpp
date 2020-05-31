#ifndef SUBTRACT_HPP_
#define SUBTRACT_HPP_

#include <map>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "subtracttypes.hpp"
using json = nlohmann::json;

namespace omm {

  Action subtract(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num - num = num
      string - num = string
      boolean - boolean = num
      array - num = array
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num - num = num

      std::string val(SubtractC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]));

      finalRet = Action{ "number", "", { val }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture };

    } else if ((num1.Type == "number" && num2.Type != "string") || (num2.Type == "string" && num1.Type != "number")) { //detect case string - num = string

      finalRet = subtractstrings(num1, num2, cli_params, this_vals, dir);

    } else if (num1.Type == "boolean" && num2.Type == "bolean") { //detect case boolean - boolean = num

      finalRet = subtractbools(num1, num2, cli_params, this_vals, dir);

    } else if ((num1.Type == "number" && num2.Type != "array") || (num2.Type == "array" && num1.Type != "number")) { //detect case array - num = array

      finalRet = subtractarrays(num1, num2, cli_params, this_vals, dir);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
