#ifndef MODULO_HPP_
#define MODULO_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../structs.hpp"
#include "../../values.hpp"
#include "../../../bind.h"
#include "../numeric/modulo.hpp"

namespace omm {

  Action modulo(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    /* TABLE OF TYPES:

      num % num = num
      default = falsey
    */

    Action finalRet;

    if (num1.Type == "number" && num2.Type == "number") { //detect case num % num = num

      finalRet = moduloNums(num1, num2, cli_params);

    } else { //detect default case

      //return undef
      finalRet = falseyVal;
    }

    return finalRet;
  }

}

#endif
