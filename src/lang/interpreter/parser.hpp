#ifndef PARSER_HPP_
#define PARSER_HPP_

#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <thread>
#include <windows.h>
#include <regex>
#include <exception>
#include "json.hpp"
#include "bind.h"
#include "structs.hpp"
#include "indexes.hpp"
#include "log_format.hpp"
#include "values.hpp"
#include "comparisons.hpp"
#include "similarity.hpp"

//file i/o
#include "../files/readfile.hpp"
#include "../files/writefile.h"
#include "../files/isDir.hpp"
#include "../files/isFile.hpp"

using namespace std;
using json = nlohmann::json;
using ulong = unsigned long;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn) {

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

            json parsed = parser(acts, calc_params, vars, dir, false, line, true).exp;

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

            Returner cond = parser(actions[i]["Condition"][0]["Condition"], calc_params, vars, dir, true, line, false);

            //while the alt statement should continue
            while (isTruthy(cond.exp)) {

              //going back to the first block when it reached the last block
              if (o >= actions[i]["Condition"].size()) o = 0;

              parser(actions[i]["Condition"][o]["Actions"], calc_params, vars, dir, true, line, false);

              o++;
            }
          }
          break;
        case 4: {

            //global

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json parsed = parser(acts, calc_params, vars, dir, false, line, true).exp;

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

            json _val = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            log_format(_val, calc_params, vars, dir, line, 2, "log");
          }
          break;
        case 6: {

            //print

            json _val = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            log_format(_val, calc_params, vars, dir, line, 2, "print");
          }
          break;
        case 8: {

            //expressionIndex

            json val = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            json index = indexesCalc(val["Hash_Values"], actions[i]["Indexes"], calc_params, vars, line, dir);

            if (expReturn) {
              vector<string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 9: {

            //group

            json acts = actions[i]["ExpAct"];

            Returner parsed = parser(acts, calc_params, vars, dir, false, line, false);

            json pVars = parsed.variables;

            //filter the variables that are not global
            for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
              if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                vars[o.value()["name"].get<string>()] = o.value();

            if (groupReturn) return Returner{ parsed.value, vars, parsed.exp, parsed.type };
          }
          break;
        case 10: {

            //process
                                                             /* process overloading */
            string name = actions[i]["Name"].get<string>() + to_string(actions[i]["Params"].size());

            if (name != "") {
              json acts = actions[i]["ExpAct"];

              json nVar = {
                {"type", "process"},
                {"name", name},
                {"value", actions[i]},
                {"valueActs", json::parse("[]")}
              };

              vars[name] = nVar;
            }

            if (expReturn) {
              vector<string> noRet;

              return Returner{ noRet, vars, actions[i], "expression" };
            }
          }
          break;
        case 11: {

            //# (call process)

                                                             /* process overloading */
            string name = actions[i]["Name"].get<string>() + to_string(actions[i]["Args"].size());

            Returner parsed;

            vector<string> noRet;

            Returner fparsed = Returner{ noRet, vars, falseyVal, "none" };

            parsed = fparsed;

            if (vars.find(name) == vars.end()) goto stopIndexing_processes;
            else {

              json var = vars[name]["value"];

              for (json it : actions[i]["Indexes"]) {

                json _index = parser(it, calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
                string index = _index.dump().substr(1, _index.dump().length() - 2);

                if (var["Hash_Values"].find(index) == var["Hash_Values"].end()) {
                  parsed = fparsed;
                  goto stopIndexing_processes;
                }

                var = parser(var["Hash_Values"][index], calc_params, vars, dir, false, line, true).exp;
              }

              if (var["Type"] != "process") {
                parsed = fparsed;
                goto stopIndexing_processes;
              }

              json params = var["Params"]
              , args = actions[i]["Args"];

              json sendVars = vars;

              for (int o = 0; o < params.size() || o < args.size(); o++) {

                json cur = {
                  {"type", "local"},
                  {"name", (string) params[o]},
                  {"value", parser(args[o], calc_params, vars, dir, false, line, true).exp},
                  {"valueActs", json::parse("[]")}
                };

                sendVars[(string) params[o]] = cur;
              }

              if (vars[name]["type"] == "process") {

                parsed = parser(var["ExpAct"], calc_params, sendVars, dir, true, line, false);

                json pVars = parsed.variables;

                //filter the variables that are not global
                for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                  if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                    vars[o.value()["name"].get<string>()] = o.value();

              }
            }

            stopIndexing_processes:
            if (expReturn) {

              json val = parsed.exp;

              vector<string> noRet;

              return Returner{ noRet, vars, val, "expression" };
            }
          }
          break;
        case 12: {

            //return

            vector<string> noRet;

            return Returner{ noRet, vars, parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp, "return" };
          }
          break;
        case 13: {

            //conditional

            for (int o = 0; o < actions[i]["Condition"].size(); o++) {

              json val = parser(actions[i]["Condition"][o]["Condition"], calc_params, vars, dir, false, line, true).exp;

              if (isTruthy(val)) {

                Returner parsed = parser(actions[i]["Condition"][o]["Actions"], calc_params, vars, dir, true, line, false);

                json pVars = parsed.variables;

                //filter the variables that are not global
                for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                  if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                    vars[o.value()["name"].get<string>()] = o.value();

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

            json files = actions[i]["Value"]; //get all actionized files imported

            //loop through actionized files
            for (json it : files) {

              Returner parsed = parser(it, calc_params, vars, dir, true, 0, false);

              json pVars = parsed.variables;

              //filter the variables that are not global
              for (auto& o : pVars.items())
                if (o.value()["type"] == "global" || o.value()["type"] == "process")
                  vars[o.key()] = o.value();
            }
          }
          break;
        case 15: {

            //read

            string in;

            cout << ((string) parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp["ExpStr"][0]) << " ";

            cin >> in;

            if (expReturn) {
              vector<string> retNo;

              json expRet = {
                {"Type", "string"},
                {"Name", ""},
                {"ExpStr", json::parse("[\"\'" + in + "\'\"]")},
                {"ExpAct", "[]"_json},
                {"Params", "[]"_json},
                {"Args", "[]"_json},
                {"Condition", "[]"_json},
                {"ID", 38},
                {"First", "[]"_json},
                {"Second", "[]"_json},
                {"Degree", "[]"_json},
                {"Value", "[[]]"_json},
                {"Indexes", "[[]]"_json},
                {"Index_Type", ""},
                {"Hash_Values", {
                  {"falsey", falseyVal}
                }},
                {"IsMutable", false}
              };

              return Returner{ retNo, vars, expRet, "expression" };
            }
          }
          break;
        case 16: {

            //break

            Returner ret;

            vector<string> returnNone;

            ret.value = returnNone;
            ret.variables = vars;
            ret.exp = "{}"_json;
            ret.type = "break";

            return ret;
          }
          break;
        case 17: {

            //skip

            Returner ret;

            vector<string> returnNone;

            ret.value = returnNone;
            ret.variables = vars;
            ret.exp = "{}"_json;
            ret.type = "skip";

            return ret;
          }
          break;
        case 18: {

            //eval

            string code = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp.get<string>();

            char* sendDir = const_cast<char*>(&dir[0]);

            char* len = CLex(&code[0], sendDir, "eval");

            char* __acts = Cactions(len, sendDir, "eval");

            string _acts(__acts);

            json acts = json::parse(_acts);

            Returner parsed = parser(acts, calc_params, vars, dir, false, line, false);

            if (expReturn) {

              vector<string> returnNone;

              return Returner{ returnNone, vars, parsed.exp, "expression" };
            }
          }
          break;
        case 19: {

            //typeof

            Returner parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true);

            json exp = parsed.exp;
            json stringval = strPlaceholder;

            stringval["ExpStr"] = json::parse("[\"" + exp["Type"].get<string>() + "\"]");

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, stringval, "expression" };
          }
          break;
        case 21: {

            //loop

            json cond = actions[i]["Condition"][0]["Condition"]
            , acts = actions[i]["Condition"][0]["Actions"];

            Returner parsed;

            json condP = parser(cond, calc_params, vars, dir, false, line, true).exp;

            while (isTruthy(condP)) {

              parsed = parser(acts, calc_params, vars, dir, true, line, false);

              json pVars = parsed.variables;

              //filter the variables that are not global
              for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                  vars[o.value()["name"].get<string>()] = o.value();

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") continue;
              if (parsed.type == "break") break;

              condP = parser(cond, calc_params, vars, dir, false, line, true).exp;
            }

          }
          break;
        case 22: {

            //hash

            if (expReturn) {

              vector<string> returnNone;

              bool isMutable = actions[i]["IsMutable"].get<bool>();

              json val = actions[i];

              if (!isMutable) {

                for (auto& it : actions[i]["Hash_Values"].items())
                  val["Hash_Values"][it.key()] = json::parse("[" + parser(it.value(), calc_params, vars, dir, false, line, true).exp.dump() + "]");
              }

              return Returner{ returnNone, vars, val, "expression" };
            }
          }
          break;
        case 23: {

            //hashIndex

            json val = actions[i]["Hash_Values"];

            json index = indexesCalc(val, actions[i]["Indexes"], calc_params, vars, line, dir);

            if (expReturn) {
              vector<string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 24: {

            //array

            if (expReturn) {
              vector<string> returnNone;

              bool isMutable = actions[i]["IsMutable"].get<bool>();

              json val = actions[i];

              if (!isMutable) {

                char* index = "0";

                for (json o : actions[i]["Hash_Values"]) {

                  if (val["Hash_Values"].find(index) == val["Hash_Values"].end()) {
                    index = AddC(index, "1");
                    continue;
                  }

                  val["Hash_Values"][index] = json::parse("[" + parser(actions[i]["Hash_Values"][index], calc_params, vars, dir, false, line, true).exp.dump() + "]");
                  index = AddC(index, "1");
                }
              }

              return Returner{ returnNone, vars, val, "expression"};
            }
          }
          break;
        case 25: {

            //arrayIndex

            json val = actions[i]["Hash_Values"]
            , index = indexesCalc(val, actions[i]["Indexes"], calc_params, vars, line, dir);

            if (expReturn) {
              vector<string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 26: {

            //ascii

            json parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            vector<string> returnNone;

            if (parsed["Type"] != "string" && expReturn) return Returner{ returnNone, vars, falseyVal, "expression" };
            else {
              string val = parsed["ExpStr"][0].get<string>().substr(1, parsed["ExpStr"][0].get<string>().length() - 2);
              int first = (int) val[0];

              if (expReturn) {

                json ascVal = {
                  {"Type", "number"},
                  {"Name", ""},
                  {"ExpStr", json::parse("[\"" + to_string(first) + "\"]")},
                  {"ExpAct", "[]"_json},
                  {"Params", "[]"_json},
                  {"Args", "[]"_json},
                  {"Condition", "[]"_json},
                  {"ID", 39},
                  {"First", "[]"_json},
                  {"Second", "[]"_json},
                  {"Degree", "[]"_json},
                  {"Value", "[[]]"_json},
                  {"Indexes", "[[]]"_json},
                  {"Index_Type", ""},
                  {"Hash_Values", {
                    {"falsey", falseyVal}
                  }},
                  {"IsMutable", false}
                };

                return Returner{returnNone, vars, ascVal, "expression"};
              }
            }
          }
          break;
        case 28: {

            //let

            string name = actions[i]["Name"];

            json acts = actions[i]["ExpAct"];

            json parsed = parser(acts, calc_params, vars, dir, false, line, true).exp;

            if (parsed.size() == 0) {
              cout << "There Was An Unidentified Error On Line " << line << endl;
              Kill();
            }

            json nVar;

            json var = vars[name];
            vector<string> indexes;

            for (json it : actions[i]["Indexes"]) {
              json varP = parser(it, calc_params, vars, dir, false, line, true).exp["ExpStr"][0];

              if (var["value"]["Hash_Values"].find(varP.get<string>()) == var["value"]["Hash_Values"].end()) var = {
                  {"type", "local"},
                  {"name", var["name"].get<string>() + varP.get<string>()},
                  {"value", {
                    {varP.get<string>(), {
                      {"falsey", falseyVal}
                    }}
                  }},
                  {"valueActs", json::parse("[]")}
                };
              else var = {
                {"type", "local"},
                {"name", var["name"].get<string>() + varP.get<string>()},
                {"value", var["value"]["Hash_Values"][varP.get<string>()]},
                {"valueActs", json::parse("[]")}
              };

              indexes.push_back(varP.get<string>());
            }

            if (var.find("type") != var.end())
              nVar = {
                {"type", vars[name]["type"]},
                {"name", name},
                {"value", parsed},
                {"valueActs", json::parse("[]")}
              };
            else
              nVar = {
                {"type", "local"},
                {"name", name},
                {"value", parsed},
                {"valueActs", json::parse("[]")}
              };

            if (actions[i]["Indexes"].size() == 0) vars[name] = nVar;
            else {
               json myObj;
               auto ref = std::ref(vars[name]["value"]["Hash_Values"]);

               for (string i : indexes) ref = ref.get()[i];

               ref.get() = json::parse("[" + nVar["value"].dump() + "]");
            }
          }
          break;
        case 32: {

            //add

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Add(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 33: {

            //subtract

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Subtract(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 34: {

            //multiply

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Multiply(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 35: {

            //divide

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Division(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 36: {

            //exponentiate

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Exponentiate(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 37: {

            //modulo

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump(2)
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump(2);

            string _val(Modulo(
              &first[0],
              &second[0],
              &calc_params.dump()[0],
              line
            ));

            json val = json::parse(_val);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 38: {

            //string

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, actions[i], "expression" };
          }
          break;
        case 39: {

            //number

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, actions[i], "expression" };
          }
          break;
        case 40: {

            //boolean

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, actions[i], "expression" };
          }
          break;
        case 41: {

            //falsey

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, actions[i], "expression" };
          }
          break;
        case 42: {

            //none

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, actions[i], "expression" };
          }
          break;
        case 43: {

            //variable

            json val;

            if (vars.find(actions[i]["Name"].get<string>()) == vars.end()) val = falseyVal;
            else {

              json var = vars[actions[i]["Name"].get<string>()]["value"];

              bool varIsMutable = var["IsMutable"].get<bool>()
              , actIsMutable = actions[i]["IsMutable"].get<bool>()
              , isMutable = varIsMutable ^ actIsMutable;

              var["IsMutable"] = isMutable;

              val = parser(json::parse("[" + var.dump() + "]"), calc_params, vars, dir, false, line, true).exp;
            }

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, val, "expression" };
          }
          break;
        case 44: {

            //type

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, json::parse(
              R"(
                {"Args":[],"Condition":[],"Degree":[],"ExpAct":[],"ExpStr":["type"],"First":[],"Hash_Values":{},"ID":44,"Index_Type":"","Indexes":[],"Name":"","Params":[],"Second":[],"Type":"type","Value":[],"ValueType":[]}
              )"
            ), "expression" };
          }
          break;
        case 46: {

            //variableIndex

            Returner parsedVal = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true);

            json index = indexesCalc(parsedVal.exp["Hash_Values"], actions[i]["Indexes"], calc_params, vars, line, dir);

            if (expReturn) {
              vector<string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 47: {

            //equals

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = equals(
              first,
              second,
              calc_params,
              vars,
              dir,
              line
            );

            if (first["Type"] != second["Type"]) val = falseRet;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 48: {

            //notEqual

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = equals(
              first,
              second,
              calc_params,
              vars,
              dir,
              line
            );

            val = val["ExpStr"][0] == "true" ? falseRet : trueRet;
            if (first["Type"] != second["Type"]) val = trueRet;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 49: {

            //greater

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = isGreater(
              first,
              second,
              calc_params,
              line
            );

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 50: {

            //less

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = isLess(
              first,
              second,
              calc_params,
              line
            );

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 51: {

            //greaterOrEqual

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = isLess(
              first,
              second,
              calc_params,
              line
            );

            val = val["ExpStr"][0] == "true" ? falseRet : trueRet;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 52: {

            //lessOrEqual

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json val = isGreater(
              first,
              second,
              calc_params,
              line
            );

            val = val["ExpStr"][0] == "true" ? falseRet : trueRet;
            if (first["Type"] != second["Type"]) val = trueRet;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 53: {

            //not

            json val = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp
            , expstr = val["ExpStr"][0]
            , retval;

            if (expstr == "false" || val["Type"] == "falsey") retval = trueRet;
            else retval = falseRet;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = retval;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 54: {

            //similar

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json retval;

            if (actions[i]["Degree"].size() == 0) retval = similarity(first, second, zero, calc_params, vars, dir, line);
            else retval = similarity(first, second, parser(actions[i]["Degree"], calc_params, vars, dir, false, line, true).exp, calc_params, vars, dir, line);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = retval;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 55: {

            //strictSimilar

            json first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp;

            json retval;

            if (actions[i]["Degree"].size() == 0) retval = strictSimilarity(first, second, zero, calc_params, vars, dir, line);
            else retval = strictSimilarity(first, second, parser(actions[i]["Degree"], calc_params, vars, dir, false, line, true).exp, calc_params, vars, dir, line);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = retval;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 56:  {

            //@ (call thread)

            string name = actions[i]["Name"].get<string>() + to_string(actions[i]["Args"].size());

            Returner parsed;

            vector<string> noRet;

            Returner fparsed = Returner{ noRet, vars, falseyVal, "none" };

            parsed = fparsed;

            if (vars.find(name) == vars.end()) goto stopIndexing_threads;
            else {

              json var = vars[name]["value"];

              for (json it : actions[i]["Indexes"]) {

                json _index = parser(it, calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
                string index = _index.dump().substr(1, _index.dump().length() - 2);

                if (var["Hash_Values"].find(index) == var["Hash_Values"].end()) {
                  parsed = fparsed;
                  goto stopIndexing_threads;
                }

                var = parser(var["Hash_Values"][index], calc_params, vars, dir, false, line, true).exp;
              }

              if (var["Type"] != "process") {
                parsed = fparsed;
                goto stopIndexing_threads;
              }

              json params = var["Params"]
              , args = actions[i]["Args"];

              json sendVars = vars;

              for (int o = 0; o < params.size() || o < args.size(); o++) {

                json cur = {
                  {"type", "local"},
                  {"name", (string) params[o]},
                  {"value", parser(args[o], calc_params, vars, dir, false, line, true).exp},
                  {"valueActs", json::parse("[]")}
                };

                sendVars[(string) params[o]] = cur;
              }

              if (vars[name]["type"] == "process") {
                thread _(parser, var["ExpAct"], calc_params, sendVars, dir, true, line, false);

                _.join();
              }
            }

            stopIndexing_threads:
            if (expReturn) {

              json val = parsed.exp;

              vector<string> noRet;

              return Returner{ noRet, vars, val, "expression" };
            }
          }
          break;
        case 57: {

            //wait

            json amt = parser(actions[i]["ExpAct"], &calc_params.dump()[0], vars, dir, false, line, true).exp;

            if (IsLessC(&(amt["ExpStr"][0].get<string>())[0], "4294967296")) Sleep((ulong) atoi(&(amt["ExpStr"][0].get<string>())[0]));
            else {
              for (char* i = "0"; (bool) IsLessC(i, &(amt["ExpStr"][0].get<string>())[0]); i = AddStrings(i, "4294967296", &calc_params.dump()[0], line)) {

                char* subtracted = SubtractStrings(&(amt["ExpStr"][0].get<string>())[0], i, &calc_params.dump()[0], line);

                if (IsLessC(
                  subtracted,
                  "4294967296"
                )) Sleep((ulong) atoi(subtracted));
                else Sleep((ulong) 4294967296);
              }
            }
          }
          break;
        case 58: {

            //cast

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              json cur = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;
              cur["Type"] = actions[i]["Name"];
              cur["Name"] = actions[i]["ExpStr"][0];

              ret.exp = retNo;
              ret.variables = vars;
              ret.exp = cur;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 59: {

            //each

            json putterVars = actions[i]["ExpStr"];
            string var1 = putterVars[0]
            , var2 = putterVars[1];

            //parse the iterator value
            json iterator = parser(actions[i]["First"] /* actions[i]["First"] is where the iterator is stored */, calc_params, vars, dir, false, line, true).exp["Hash_Values"];

            iterator.erase("falsey");

            for (auto& it : iterator.items()) {
              json sendVars = vars;

              json key = {
                {"Type", "string"},
                {"Name", ""},
                {"ExpStr", json::parse("[\"" + it.key() + "\"]")},
                {"ExpAct", "[]"_json},
                {"Params", "[]"_json},
                {"Args", "[]"_json},
                {"Condition", "[]"_json},
                {"ID", 38},
                {"First", "[]"_json},
                {"Second", "[]"_json},
                {"Degree", "[]"_json},
                {"Value", "[[]]"_json},
                {"Indexes", "[[]]"_json},
                {"Index_Type", ""},
                {"Hash_Values", {
                  {"falsey", falseyVal}
                }},
                {"IsMutable", false}
              };

              sendVars[var1] = {
                {"type", "local"},
                {"name", var1},
                {"value", key},
                {"valueActs", json::parse("[]")}
              };
              sendVars[var2] = {
                {"type", "local"},
                {"name", var2},
                {"value", parser(it.value(), calc_params, vars, dir, false, line, true).exp},
                {"valueActs", json::parse("[]")}
              };

              Returner parsed = parser(actions[i]["ExpAct"], calc_params, sendVars, dir, true, line, false);

              json pVars = parsed.variables;

              //filter the variables that are not global
              for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                if (o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end())
                  vars[o.value()["name"].get<string>()] = o.value();

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") continue;
              if (parsed.type == "break") break;
            }
          }
          break;

        //all of the omm cprocs
        case 60: {

          //files.read

          //written as files.read(dir)

          string filename = parser(actions[i]["Args"][0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

          smatch match;

          //see if the filename is absolute
          regex pat("^[a-zA-Z]:");
          bool isOnDrive = regex_search(filename, match, pat);

          string nDir = isOnDrive ? "" : dir;

          if (!isFile(nDir + filename) && expReturn) {

            Returner ret;

            vector<string> retNo;

            ret.exp = retNo;
            ret.variables = vars;
            ret.exp = falseyVal;
            ret.type = "expression";

            return ret;

          } else {
            string content = readfile(&(nDir + filename)[0]);

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              json contentJ = strPlaceholder;

              contentJ["ExpStr"] = {content};

              for (ulong i = 0; i < content.length(); i++) {
                json curChar = strPlaceholder;

                curChar["ExpStr"] = {
                  to_string(content[i])
                };

                contentJ["Hash_Values"][to_string(i)] = curChar;
              }

              ret.exp = retNo;
              ret.variables = vars;
              ret.exp = contentJ;
              ret.type = "expression";

              return ret;
            }
          }

          break;
        }
        case 61: {

          //files.write

          //written as files.write(dir, content)

          string filename = parser(actions[i]["Args"][0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();
          json content = parser(actions[i]["Args"][1], calc_params, vars, dir, false, line, true).exp;

          string contentstr = content["ExpStr"][0].get<string>();

          smatch match;

          //see if the filename is absolute
          regex pat("^[a-zA-Z]:");
          bool isOnDrive = regex_search(filename, match, pat);

          string nDir = isOnDrive ? "" : dir;

          writefile(&(nDir + filename)[0], &contentstr[0]);

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            ret.exp = retNo;
            ret.variables = vars;
            ret.exp = content;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 62: {

          //files.exists

          //written as file.exists(dir)

          string filename = parser(actions[i]["Args"][0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

          smatch match;

          //see if the filename is absolute
          regex pat("^[a-zA-Z]:");
          bool isOnDrive = regex_search(filename, match, pat);

          string nDir = isOnDrive ? "" : dir;

          //if it is not a directory and not a file, it does not exist
          bool exists = !(!isDir(nDir + filename) && !isFile(nDir + filename));

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            ret.exp = retNo;
            ret.variables = vars;
            ret.exp = exists ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 63: {

          //files.isFile

          //written as file.isFile(dir)

          string filename = parser(actions[i]["Args"][0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

          smatch match;

          //see if the filename is absolute
          regex pat("^[a-zA-Z]:");
          bool isOnDrive = regex_search(filename, match, pat);

          string nDir = isOnDrive ? "" : dir;

          bool isFileVal = isFile(nDir + filename);

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            ret.exp = retNo;
            ret.variables = vars;
            ret.exp = isFileVal ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 64: {

          //files.isDir

          //written as file.isDir(dir)

          string filename = parser(actions[i]["Args"][0], calc_params, vars, dir, false, line, true).exp["ExpStr"][0].get<string>();

          smatch match;

          //see if the filename is absolute
          regex pat("^[a-zA-Z]:");
          bool isOnDrive = regex_search(filename, match, pat);

          string nDir = isOnDrive ? "" : dir;

          bool isDirVal = isDir(nDir + filename);

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            ret.exp = retNo;
            ret.variables = vars;
            ret.exp = isDirVal ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        //////////////////////////

        case 65: {

          //kill_thread

          terminate();

          break;
        }
        case 66: {

          //kill

          Kill();

          break;
        }

        //assignment operators
        case 4343: {

          //++

          string name = actions[i]["Name"];

          json nVar;

          if (vars[name]["type"] != "dynamic") {

            if (vars[name].find("value") != vars[name].end()) {

              json _val = vars[name]["value"];

              char* _added = Add(&(_val.dump())[0], &val1.dump()[0], &calc_params.dump()[0], line);
              string added(_added);

              nVar = {
                {"type", vars[name]["type"]},
                {"name", name},
                {"value", json::parse(added)},
                {"valueActs", json::parse("[]")}
              };
            } else nVar = {
                {"type", "local"},
                {"name", name},
                {"value", val1},
                {"valueActs", json::parse("[]")}
              };
          }

          vars[name] = nVar;

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val1;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 4545: {

            //--

            string name = actions[i]["Name"];

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Subtract(&(_val.dump())[0], &val1.dump()[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", valn1},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val1;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 4361: {

            //+=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Add(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 4561: {

            //-=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Subtract(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 4261: {

            //*=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Multiply(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 4761: {

            ///=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Division(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 9461: {

            //^=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Exponentiate(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 3761: {

            //%=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp;

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                json _val = vars[name]["value"];

                char* _added = Modulo(&(_val.dump())[0], &(_inc.dump())[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse(added)},
                  {"valueActs", json::parse("[]")}
                };
              } else nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", _inc},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              vector<string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
      }
    } catch (int e) {
      cout << "There Was An Unidentified Error On Line " << line << endl;
      Kill();
    }
  }

  vector<string> returnNone;

  Returner ret;

  ret.value = returnNone;
  ret.variables = vars;
  ret.exp = falseyVal;
  ret.type = "none";

  return ret;
}

#endif
