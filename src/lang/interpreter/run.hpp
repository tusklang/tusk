#ifndef RUN_HPP_
#define RUN_HPP_

#include <vector>
#include <string>
#include "parser.hpp"
#include "json.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

void run(char *acts, char *calc_params, char *dir) {

  try {

    json actions = json::parse(acts)
    , cp = json::parse(calc_params)
    , vars = json::parse("{}");

    parser(actions, cp, vars, dir, /*group return*/ false, /*line*/ 1, /* expression return */ false);
  } catch (int e) {
    cout << "There Was An Unidentified Error" << endl;
    Kill();
  }
}

#endif