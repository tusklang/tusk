#ifndef OMM_RUN_HPP_
#define OMM_RUN_HPP_

//the max digit number (if digit surpasses this, overflow it)

/*

  The reason DigitSize must be 1 is because of zeros
  example:
    DigitSize = 2
    Number1 = 1000
    Number2 = 100

  1000 -> [10, 00] -> [10, 0]
  100 -> [10, 0] -> [10, 0]

  thus 1000 == 100, but that is not really true
*/

#define DigitSize 1 //if you change here, change in numconv.go

#define OMM_MAX_DIGIT (std::pow(10, DigitSize)) //the actual max digit + 1
#define OMM_MIN_DIGIT (-1 * OMM_MAX_DIGIT) //the actual min digit - 1

#include <iostream>
#include <vector>
#include <string>
#include "parser.hpp"
#include "structs.hpp"
#include "decode_json.hpp"
#include "json.hpp"
#include "values.hpp"
#include "ommtypes.hpp"
#include "cprocs.hpp"

#include "operations/numeric/modulo.hpp"
#include "operations/numeric/normalize.hpp"
using json = nlohmann::json;

namespace omm {

  void run(char* actions, char* cli_params, char* dir, int argc, char ** argv) {

    const json cpJ = json::parse(std::string(cli_params));

    //convert the json to a vector of actions
    std::vector<Action> acts = DecodeJSON::vector(json::parse(std::string(actions))); //maybe in the future remove json entirely from this, but I dont know how

    std::map<std::string, Variable> vars;

    Action dirnameAct = { "string", "", { string(dir) }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };

    vars["$dirname"] = Variable{
      "global",
      "$dirname",
      { dirnameAct },
      [](Action v, json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) -> Returner { return Returner{}; }
    };

    Variable omm_args = Variable{
      "global",
      "$argv",
      { arrayVal },
      [](Action v, json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) -> Returner { return Returner{}; }
    };

    for (int i = 0; i < argc; ++i)
      omm_args.value[0].Hash_Values[to_string(i)] = { ommtypes::to_string(string(argv[i])) };

    vars["$argv"] = omm_args;

    //put all of the cprocs into the vars
    for (std::pair<std::string, std::function<Returner(
      Action v,
      json cli_params,
      std::map<std::string, Variable> vars,
      std::deque<std::map<std::string, std::vector<Action>>> this_vals,
      std::string dir
    )>> it : cprocs)
      vars["$" + it.first] = Variable{
        "cproc",
        "$" + it.first,
        {},
        it.second
      };

    parser(acts, cpJ, vars, /*group return*/ false, /* expression return */ false, {}, std::string(dir));
  }

}

#endif
