#ifndef ADD_HPP_
#define ADD_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "addtypes.hpp"
#include "../numeric/add.hpp"
using json = nlohmann::json;

namespace omm {

  Action add(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      string + (* - array - none - hash) = string
      array + (* - none) = array
      num + num = num
      hash + hash = hash
      boolean + boolean = boolean
      num + boolean = num
      default = falsey
    */

    Action finalRet;

    if ((num1.Type == "string" || num2.Type == "string") && ((num1.Type != "array" && num2.Type != "array") && (num1.Type != "none" && num2.Type != "none") && (num1.Type != "hash" && num2.Type != "hash"))) { //detect case string + (* - array - none - hash) = string

      finalRet = addstrings(num1, num2, cli_params, this_vals, dir);

    } else if ((num1.Type == "array" || num2.Type == "array") && (num1.Type != "none" && num2.Type != "none")) { //detect case array + (* - none) = array

      finalRet = addarrays(num1, num2, cli_params, this_vals, dir);

    } else if ((num1.Type == "number" || num2.Type == "number") && (num1.Type == "number" || num2.Type == "number")) { //detect case num + num = num

      finalRet = addNums(num1, num2, cli_params);

    } else if ((num1.Type == "hash" || num2.Type == "hash") && (num1.Type == "hash" || num2.Type == "hash")) { //detect case hash + hash = hash

      finalRet = addhashes(num1, num2, cli_params, this_vals, dir);

    } else if ((num1.Type == "boolean" || num2.Type == "boolean") && (num1.Type == "boolean" || num2.Type == "boolean")) { //detect case boolean + boolean = boolean

      finalRet = addbools(num1, num2, cli_params, this_vals, dir);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
