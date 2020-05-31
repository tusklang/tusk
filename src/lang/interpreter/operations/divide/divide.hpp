#ifndef DIVIDE_HPP_
#define DIVIDE_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
using json = nlohmann::json;

namespace omm {

  Action divide(Action num1, Action num2, json cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num / num = num
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num / num = num

      std::string val(DivisionC(&num1.ExpStr[0][0], &num2.ExpStr[0][0], &cli_params.dump()[0]));

      finalRet = Action{ "number", "", { val }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyFuture };

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
