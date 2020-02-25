#include <iostream>
#include <vector>
#include "json.hpp"
#include "math.hpp"
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

Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line) {

  json expStr;

  for (int i = 0; i < actions.size(); i++) {

    int cur = actions[i]["ID"];

    try {
      switch (cur) {
        case 0: //newline
          line++;
          break;
        case 1: {
            //local

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json parsed = parser(acts, calc_params, vars, dir, false, line).exp;

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
          }
          break;
        case 2: {
            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json nVar = {
              {"type", "local"},
              {"name", name},
              {"value", acts},
              {"valueActs", json::parse("[]")}
            };
            vars[name] = nVar;
          }
          break;
        case 3: {
            //alt

            int o = 0;

            struct Returner cond = parser(actions[i]["Condition"][0]["Condition"], calc_params, vars, dir, true, line);

            //while the alt statement should continue
            while (cond.exp[0][0] != "false" && cond.exp[0][0] != "undefined" && cond.exp[0][0] != "null") {

              //going back to the first block when it reached the last block
              if (o >= actions[i]["Condition"].size()) o = 0;

              parser(actions[i]["Condition"][o]["Actions"], calc_params, vars, dir, true, line);

              o++;
            }
          }
          break;
        case 5:
          cout << parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0].dump() << endl;
          break;
        case 6:
          cout << parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0].dump();
          break;
        case 7: {
            //expression

            string expStr = actions[i]["ExpStr"].dump();

            json nExp = json::parse("[" + expStr + "]");

            cout << nExp << endl;

            json calculated = math(nExp, calc_params, vars, dir, line);

            expStr = calculated;
          }
          break;
      }
    } catch (int e) {
      cout << "There Was An Unidentified Error On Line " << line << endl;
      Kill();
    }
  }

  vector<string> returnNone;

  return Returner{ returnNone, vars, math(expStr, calc_params, vars, dir, line), "none" };
}
