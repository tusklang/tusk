#ifndef RUN_HPP_
#define RUN_HPP_

#include <iostream>
#include <vector>
#include <string>
#include "parser.hpp"
#include "structs.hpp"
#include "decode_json.hpp"
#include "json.hpp"
#include "values.hpp"
#include "ommtypes.hpp"
using namespace std;
using json = nlohmann::json;

void run(char* actions, char* cli_params, char* dir, int argc, char ** argv) {

  const json cpJ = json::parse(string(cli_params));

  //convert the json to a vector of actions
  vector<Action> acts = DecodeJSON::vector(json::parse(string(actions)));

  map<string, Variable> vars;

  vars["$dirname"] = Variable{
    "global",
    "$dirname",
    {
      { "string", "", { string(dir) }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" }
    }
  };

  Variable omm_args = Variable{
    "global",
    "$argv",
    { arrayVal }
  };

  for (int i = 0; i < argc; ++i)
    omm_args.value[0].Hash_Values[to_string(i)] = { ommtypes::to_string(string(argv[i])) };

  vars["$argv"] = omm_args;

  parser(acts, cpJ, vars, /*group return*/ false, /* expression return */ false, {}, string(dir));
}

#endif
