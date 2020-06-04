#ifndef OMM_INDEXES_HPP_
#define OMM_INDEXES_HPP_

#include <map>
#include <vector>
#include <deque>
#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

namespace omm {

  Returner parser(const std::vector<Action> actions, const json cli_params, std::map<std::string, Variable> vars, const bool groupReturn, const bool expReturn, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir);

  Action indexesCalc(std::map<string, std::vector<Action>> val, std::vector<std::vector<Action>> indexes, json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    if (indexes.size() == 0) return falseyVal;

    std::string index = parser(indexes[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

    if (val.find(index) == val.end()) { //the second boolean is determining if the value is public

      if (val.find("falsey") == val.end()) return falseyVal;
      else return val["falsey"][0];

    } else {

      //if it is a private variable return undef
      if (std::islower(index[0]) && val[index][0].Access != "public")
        if (val.find("falsey") == val.end()) return falseyVal;
        else return val["falsey"][0];

      Action expVal = parser(val[index], cli_params, vars, false, true, this_vals, dir).exp;

      indexes.erase(indexes.begin());

      if (indexes.size() == 0) return expVal;

      return indexesCalc(expVal.Hash_Values, indexes, cli_params, vars, this_vals, dir);
    }

  }

}

#endif
