#include "parser.hpp"
#include "bridge.hpp"
#include "json.hpp"
#include <iostream>
using namespace std;
using json = nlohmann::json;

int bridge(char[] val, char[] calc_params, char[] vars, char[] dir) {
  cout << val;

  json acts = val;
  parser(acts, "{}", "{}", "");
}
