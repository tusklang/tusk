#ifndef EXPONENTIATETYPES_HPP_
#define EXPONENTIATETYPES_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../CliParams.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
using CliParams = nlohmann::CliParams;

namespace omm {

  Action exponentiatenumbers(Action num1, Action num2, CliParams cli_params, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    //for now just return undef
    return falseyVal;
  }

}

#endif
