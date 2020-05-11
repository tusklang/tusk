#ifndef MODULO_HPP_
#define MODULO_HPP_

#include <map>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "../../../bind.h"
using namespace std;
using json = nlohmann::json;

Action modulo(Action num1, Action num2, json cli_params) {

  /* TABLE OF TYPES:

    num % num = num
    default = falsey
  */

  Action finalRet;

  if (num1.Type == "number" && num2.Type == "number") { //detect case num % num = num

    if (strcmp(ReturnInitC(&num2.ExpStr[0][0]), "0") == 0)  {
      finalRet = Action{ "number", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };
    } else if (strcmp(ReturnInitC(&num1.ExpStr[0][0]), "0") == 0)  {
      finalRet = Action{ "number", "", { "0" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };
    } else {

      char* divved_ = DivisionC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]);

      string divved(divved_);

      divved = divved.substr(0, divved.find("."));

      char* mult = MultiplyC(&divved[0], &num2.ExpStr[0][0], &cli_params.dump()[0])
      ,* remainder = SubtractC(&num1.ExpStr[0][0], mult, &cli_params.dump()[0]);

      finalRet = Action{ "number", "", { remainder }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false };
    }

  } else { //detect default case

    //return undef
    finalRet = falseyVal;
  }

  return finalRet;
}

#endif
