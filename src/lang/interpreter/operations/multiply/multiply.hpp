#ifndef MULTIPLY_HPP_
#define MULTIPLY_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "multiplytypes.hpp"
using namespace std;
using json = nlohmann::json;

Action multiply(Action num1, Action num2, json cli_params, deque<map<string, vector<Action>>> this_vals) {

  /* TABLE OF TYPES:

    num * num = num
    string * num = string
    array * array = array
    default = falsey
  */

  Action finalRet;

  if (num1.Type == "number" && num2.Type == "number") { //detect case num * num = num

    string val(MultiplyC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]));

    finalRet = Action{ "number", "", { val }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };

  } else if ((num1.Type == "string" && num2.Type == "number") || (num1.Type == "number" && num2.Type == "string")) { //detect case string * num = string

    finalRet = multiplystrings(num1, num2, cli_params, this_vals);

  } else if (num1.Type == "array" && num2.Type == "array") { //detect case array * array = array

    finalRet = multiplyarrays(num1, num2, cli_params, this_vals);

  } else { //detect default case

    //return undef
    finalRet = falseyVal;
  }

  return finalRet;
}

#endif
