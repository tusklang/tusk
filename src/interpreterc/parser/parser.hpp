#include <iostream>
#include <vector>
#include <string>
#include "json.hpp"
#include "bind.h"
#include "indexes.hpp"
#include "math.hpp"
#include "structs.h"
using namespace std;
using json = nlohmann::json;

json math(json exp, const json calc_params, json vars, const string dir, int line);

Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line) {

  //empty expStr
  json expStr = "[]"_json;

  //loop through every action
  for (int i = 0; i < actions.size(); i++) {

    //get current action id
    int cur = actions[i]["ID"];

    try {
      switch (cur) {
        case 0:

          //newline

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

            //dynamic

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json nVar = {
              {"type", "dynamic"},
              {"name", name},
              {"value", json::parse("[]")},
              {"valueActs", acts}
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
        case 4: {

            //global

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json parsed = parser(acts, calc_params, vars, dir, false, line).exp;

            if (parsed.size() == 0) {
              cout << "There Was An Unidentified Error On Line " << line << endl;
              Kill();
            }

            json nVar = {
              {"type", "global"},
              {"name", name},
              {"value", parsed},
              {"valueActs", json::parse("[]")}
            };

            vars[name] = nVar;
          }
          break;
        case 5: {

            //log

            string val = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0].dump();

            val = val.substr(1);
            val.pop_back();

            cout << val << endl;
          }
          break;
        case 6: {

            //print

            string val = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0].dump();

            val = val.substr(1);
            val.pop_back();

            cout << val;
          }
          break;
        case 7: {

            //expression

            string expStr_ = actions[i]["ExpStr"].dump();

            json nExp = json::parse("[" + expStr_ + "]");

            json calculated = math(nExp, calc_params, vars, dir, line);

            expStr.push_back(calculated[0]);
          }
          break;
        case 8: {

            //expressionIndex

            string expStr_ = actions[i]["ExpStr"].dump();

            json nExp = json::parse("[" + expStr_ + "]");

            json calculated = math(nExp, calc_params, vars, dir, line);

            json index = indexesCalc(calculated, actions[i]["Indexes"], calc_params, line);

            expStr.push_back(index);
          }
          break;
        case 9: {

            //group

            json acts = actions[i]["ExpAct"];

            Returner parsed = parser(acts, calc_params, vars, dir, false, line);

            //filter the variables that are not global
            for (auto o = parsed.variables.begin(); o != parsed.variables.end(); o++)
              if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                parsed.variables.erase(o);

            vars.insert(parsed.variables.begin(), parsed.variables.end());

            if (groupReturn) {

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") return Returner{ parsed.value, vars, parsed.exp, "skip" };
              if (parsed.type == "break") return Returner{ parsed.value, vars, parsed.exp, "break" };

              return Returner{ parsed.value, vars, parsed.exp, parsed.type };
            }
          }
          break;
        case 10: {

            //process

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json nVar = {
              {"type", "process"},
              {"name", name},
              {"value", json::parse("[]")},
              {"valueActs", acts},
              {"params", actions[i]["Params"]}
            };
            vars[name] = nVar;
          }
          break;
        case 11: {

            //# (call process)

            string name = actions[i]["Name"];

            json var = vars[name];

            json params = var["params"]
            , args = actions[i]["Args"];

            json nParams = "{}"_json;

            for (int o = 0; o < params.size() || o < args.size(); o++) {

              json cur = {
                {"type", "local"},
                {"name", (string) params[o]},
                {"value", parser(json::parse("[" + args[o].dump() + "]"), calc_params, vars, dir, false, line).exp},
                {"valueActs", json::parse("[]")}
              };

              nParams[(string) params[o]] = cur;
            }

            json sendVars = vars;

            for (auto it = nParams.begin(); it != nParams.end(); it++)
              sendVars[it.key()] = it.value();

            Returner parsed = parser(var["valueActs"], calc_params, sendVars, dir, true, line);

            //filter the variables that are not global
            for (auto o = parsed.variables.begin(); o != parsed.variables.end(); o++)
              if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                parsed.variables.erase(o);

            vars.insert(parsed.variables.begin(), parsed.variables.end());

            expStr.push_back((json) parsed.value);
          }
          break;
        case 12:

          //return

          return Returner{ parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0], vars, expStr, "return" };
          break;
        case 13: {

            //conditional

            for (int o = 0; o < actions[i]["Condition"].size(); o++) {

              string val = (string) parser(actions[i]["Condition"][o]["Condition"], calc_params, vars, dir, false, line).exp[0][0];

              if (val != "false" && val != "undefined" && val != "null") {

                Returner parsed = parser(actions[i]["Condition"][o]["Actions"], calc_params, vars, dir, true, line);

                //filter the variables that are not global
                for (auto o = parsed.variables.begin(); o != parsed.variables.end(); o++)
                  if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                    parsed.variables.erase(o);

                vars.insert(parsed.variables.begin(), parsed.variables.end());

                if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
                if (parsed.type == "skip") return Returner{ parsed.value, vars, parsed.exp, "skip" };
                if (parsed.type == "break") return Returner{ parsed.value, vars, parsed.exp, "break" };

                break;
              }

            }
          }
          break;
        case 14: {

            //import

            string fileName = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0];

            if (fileName.rfind("\'", 0) == 0 || fileName.rfind("\"", 0) == 0 || fileName.rfind("`", 0) == 0) fileName = fileName.substr(1, fileName.length() - 2);

            string readerFile = dir + fileName
            , errMsg = "Could Not Find File: " + fileName;

            char* file = CReadFile(&readerFile[0], &errMsg[0], 1);

            string _acts = Cactions( CLex(file) );

            json acts = json::parse(_acts);

            Returner parsed = parser(acts, calc_params, vars, dir, false, 1);

            //filter the variables that are not global
            for (auto o = parsed.variables.begin(); o != parsed.variables.end(); o++)
              if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                parsed.variables.erase(o);

            vars.insert(parsed.variables.begin(), parsed.variables.end());
          }
          break;
        case 15: {

            //read

            string in;

            cout << ((string) parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0]) << " ";

            cin >> in;

            expStr.push_back(json::parse("[\"\'" + in + "\'\"]"));
          }
          break;
        case 16:

          //break
          return Returner{ "[]"_json, vars, expStr, "break" };
          break;
        case 17:

          //skip
          return Returner{ "[]"_json, vars, expStr, "skip" };
          break;
        case 18: {

            //eval

            string _code = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0].dump()
            , code = _code.substr(2, _code.length() - 4);

            char* codeNQ = NQReplaceC(&code[0]);

            char* lex = CLex(codeNQ);

            char* __acts = Cactions(lex);

            string _acts(__acts);

            json acts = json::parse(_acts);

            Returner parsed = parser(acts, calc_params, vars, dir, false, line);

            expStr.push_back(json::parse("[\"" + parsed.value[0] + "\"]"));
          }
          break;
        case 19: {

            //typeof

            Returner parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line);

            string exp = parsed.exp[0][0].dump().substr(1, parsed.exp[0][0].dump().length() - 2);

            char* _type = GetType(&exp[0]);

            string type(_type);

            expStr.push_back(json::parse("[\"" + type + "\"]"));
          }
          break;
        case 20: {

            //err

            Returner parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line);

            string exp = parsed.exp[0][0].dump().substr(1, parsed.exp[0][0].dump().length() - 2);

            cout << exp << "\n\nerr~" << exp << "\n^^^^ <-- Error On Line " << line << endl;

            Kill();
          }
          break;
        case 21: {

            //loop

            json cond = actions[i]["Condition"][0]["Condition"]
            , acts = actions[i]["Condition"][0]["Actions"];

            json condP = parser(cond, calc_params, vars, dir, false, line).exp[0][0];

            while (condP != "false" && condP != "undefined" && condP != "null") {

              Returner parsed = parser(acts, calc_params, vars, dir, true, line);

              //filter the variables that are not global
              for (auto o = parsed.variables.begin(); o != parsed.variables.end(); o++)
                if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                  parsed.variables.erase(o);

              vars.insert(parsed.variables.begin(), parsed.variables.end());

              condP = parser(cond, calc_params, vars, dir, false, line).exp[0][0];
            }
          }
          break;
        case 22: {

            //hash

            json expStr_ = json::parse(actions[i]["ExpStr"].dump());

            expStr.push_back(expStr_);
          }
          break;
        case 23: {

            //hashIndex

            string expStr_ = actions[i]["ExpStr"].dump();

            json nExp = json::parse("[" + expStr_ + "]");

            json calculated = math(nExp, calc_params, vars, dir, line);

            json index = indexesCalc(calculated, actions[i]["Indexes"], calc_params, line);

            expStr.push_back(index);
          }
          break;
        case 24: {

            //array

            json expStr_ = json::parse(actions[i]["ExpStr"].dump());

            expStr.push_back(expStr_);
          }
          break;
        case 25: {

            //arrayIndex

            string expStr_ = actions[i]["ExpStr"].dump();

            json nExp = json::parse("[" + expStr_ + "]");

            json calculated = math(nExp, calc_params, vars, dir, line);

            json index = indexesCalc(calculated, actions[i]["Indexes"], calc_params, line);

            expStr.push_back(index);
          }
          break;
        case 26: {

            //ascii

            string parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0].dump();

            parsed = parsed.substr(1, parsed.length() - 2);

            char first = parsed[0];

            expStr.push_back(json::parse("[\"" + to_string((int) first) + "\"]"));
          }
          break;
        case 27: {

            //parse

            string parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line).exp[0][0].dump();

            parsed = parsed.substr(1, parsed.length() - 2);

            if (!(strcmp(GetType(&parsed[0]), "string") == 0 || strcmp(GetType(&parsed[0]), "number") == 0)) {
              cout << "There Was An Error: `parse~` cannot be used on a non-string or number" << "\n\nparse~" + parsed << "\n^ <-- Error On Line " + line;
              Kill();
            }

            if (strcmp(GetType(&parsed[0]), "string") == 0)
              expStr.push_back(json::parse("[\"" + parsed.substr(1, parsed.length() - 2) + "\"]"));
            else expStr.push_back(json::parse("[\"" + parsed + "\"]"));
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
