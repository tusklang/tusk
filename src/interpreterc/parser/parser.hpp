#include <iostream>
#include <vector>
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

struct Variable {
  string                  type;
  string                  name;
  vector<string>          value;
  vector<Action>          valueActs;
};

struct Returner {
  vector<string>          value;
  json                    variables;
  vector<vector<string> > exp;
  string                  type;
};

void parser(json actions, json calc_params, json vars, string dir, bool groupReturn, int line) {
  for (int i = 0; i < actions.size(); i++) {

    int cur = actions[i]["ID"];

    switch (cur) {
      case 0:
        line++;
        break;
      case 1:

    }
  }

  cout << line << endl;
}
