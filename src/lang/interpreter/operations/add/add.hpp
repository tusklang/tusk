#ifndef ADD_HPP_
#define ADD_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "addtypes.hpp"
using namespace std;
using json = nlohmann::json;

Action add(Action num1, Action num2, json cli_params, deque<map<string, vector<Action>>> this_vals, string dir) {

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

    string val(AddC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]));

    finalRet = Action{ "number", "", { val }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };

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

#endif
