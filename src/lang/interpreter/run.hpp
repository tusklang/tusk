#ifndef RUN_HPP_
#define RUN_HPP_

#include <vector>
#include <string>
#include "parser.hpp"
#include "json.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

void run(char *acts, char *cli_params) {

  try {

    json actions = json::parse(acts)
    , cp = json::parse(cli_params)
    , vars = json::parse("{}");

    parser(actions, cp, vars, /*group return*/ false, /* expression return */ false);
  } catch (int e) {
    cout << "There Was An Unidentified Error" << endl;
    Kill();
  }
}

#endif
