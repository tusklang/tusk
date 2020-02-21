#include <vector>
#include <string>
#include "parser.hpp"
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

void run(char *acts, char *calc_params, char *dir) {
  json actions = json::parse(acts)
  , cp = json::parse(calc_params)
  , vars = json::parse("{}");

  parser(actions, cp, vars, dir, /*group return*/ false, /*line*/ 1);
}
