#ifndef EXPONENTIATE_HPP_
#define EXPONENTIATE_HPP_

#include <map>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "exponentiatetypes.hpp"
using namespace std;
using json = nlohmann::json;

Action exponentiate(Action num1, Action num2, json cli_params) {

  /* TABLE OF TYPES:

    num ^ num = num
    default = num
  */

  Action finalRet;

  if (num1.Type == "number" && num2.Type == "number") { //detect case num ^ num = num

    finalRet = exponentiatenumbers(num1, num2, cli_params);

  } else {

    //return undef
    finalRet = falseyVal;
  }

  return finalRet;
}

#endif
