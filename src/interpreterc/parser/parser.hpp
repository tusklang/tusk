#include <iostream>
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

void parser(json actions, json calc_params, string dir) {
  cout << actions.size();
}
