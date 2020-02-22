#include <iostream>
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
  vector<vector<string> > exp;
  string                  type;
};

Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line) {

  vector<vector<string>> expStr;

  for (int i = 0; i < actions.size(); i++) {

    int cur = actions[i]["ID"];

    try {
      switch (cur) {
        case 0: //newline
          line++;
          break;
        case 1: //local
          string name = actions[i]["Name"];

          json acts = actions[i]["ExpAct"];

          vector<vector<string>> parsed = parser(acts, calc_params, vars, dir, false, line).exp;

          if (parsed.size() == 0) {
            cout << "There Was An Unidentified Error On Line " << line << endl;
            Kill();
          }

          json nVar = {
            {"type", "local"},
            {"name", name},
            {"value", parsed},
            {"valueActs", json::parse("[]")}
          };

          vars[name] = nVar;
          break;
        case 7: //expression
        //link to fortran
      }
    } catch (int e) {
      cout << "There Was An Unidentified Error On Line " << line << endl;
      Kill();
    }
  }

  vector<string> returnNone;

  return Returner{ returnNone, vars, expStr, "none" };
}
