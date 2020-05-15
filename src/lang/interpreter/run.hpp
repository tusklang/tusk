#ifndef RUN_HPP_
#define RUN_HPP_

#include <iostream>
#include <vector>
#include <string>
#include "parser.hpp"
#include "structs.hpp"
#include "decode_json.hpp"
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

void run(char* actions, char* cli_params, char* dir) {

  const json cpJ = json::parse(string(cli_params));

  //convert the json to a vector of actions
  vector<Action> acts = DecodeJSON::vector(json::parse(string(actions)));

  map<string, Variable> vars;

  parser(acts, cpJ, vars, /*group return*/ false, /* expression return */ false, {}, string(dir));
}

#endif
