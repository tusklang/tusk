#ifndef STRUCTS_HPP_
#define STRUCTS_HPP_

#include <vector>
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

struct Variable {
  string                  type;
  string                  name;
  vector<string>          value;
  json                    valueActs;
};

struct Returner {
  vector<string>          value;
  json                    variables;
  json                    exp;
  string                  type;
};

#endif
