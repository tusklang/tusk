#ifndef MODULO_HPP_
#define MODULO_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
#include "../../../bind.h"
using namespace std;
using json = nlohmann::json;

Action modulo(Action num1, Action num2, json cli_params, deque<map<string, vector<Action>>> this_vals) {

  /* TABLE OF TYPES:

    num % num = num
    default = falsey
  */

  Action finalRet;

  if (num1.Type == "number" && num2.Type == "number") { //detect case num % num = num

    if (strcmp(ReturnInitC(&num2.ExpStr[0][0]), "0") == 0)  {
      finalRet = Action{ "number", "", { "undef" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };
    } else if (strcmp(ReturnInitC(&num1.ExpStr[0][0]), "0") == 0)  {
      finalRet = Action{ "number", "", { "0" }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };
    } else {

      char* divved_ = DivisionC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]);

      string divved(divved_);

      divved = divved.substr(0, divved.find("."));

      char* mult = MultiplyC(&divved[0], &num2.ExpStr[0][0], &cli_params.dump()[0])
      ,* remainder = SubtractC(&num1.ExpStr[0][0], mult, &cli_params.dump()[0]);

      finalRet = Action{ "number", "", { remainder }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };
    }

  } else { //detect default case

    //return undef
    finalRet = falseyVal;
  }

  return finalRet;
}

#endif
