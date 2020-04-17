#ifndef PARSER_HPP_
#define PARSER_HPP_

#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <thread>
#include "json.hpp"
#include "bind.h"
#include "structs.h"
#include "indexes.hpp"
#include "log_format.hpp"
#include "falsey_val.hpp"
using namespace std;
using json = nlohmann::json;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn) {

  //empty expStr
  json expStr = "[[]]"_json;

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

            cout << parsed.dump(2) << endl;

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

            struct Returner cond = parser(actions[i]["Condition"][0]["Condition"], calc_params, vars, dir, true, line, false);

            //while the alt statement should continue
            while (cond.exp[0][0] != "false" && cond.exp[0][0] != "undefined" && cond.exp[0][0] != "null") {

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

            Returner parsedVal = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true);

            json index = indexesCalc(parsedVal.exp["Hash_Values"], actions[i]["Indexes"], calc_params, vars, line, dir);

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
              if (!(o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end()))
                vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();

            if (groupReturn) return Returner{ parsed.value, vars, parsed.exp, parsed.type };
          }
          break;
        case 10: {

            //process

            string name = actions[i]["Name"];

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

            string name = actions[i]["Name"];

            Returner parsed;

            vector<string> noRet;
            json type = {
              {"Type", "falsey"},
              {"Name", ""},
              {"ExpStr", json::parse("[\"undefined\"]")},
              {"ExpAct", "[]"_json},
              {"Params", "[]"_json},
              {"Args", "[]"_json},
              {"Condition", "[]"_json},
              {"ID", 41},
              {"First", "[]"_json},
              {"Second", "[]"_json},
              {"Degree", "[]"_json},
              {"Value", "[[]]"_json},
              {"Indexes", "[[]]"_json},
              {"Index_Type", ""},
              {"Hash_Values", "{}"_json},
              {"ValueType", "[]"_json}
            }
            , fRet = {
              {"Type", "falsey"},
              {"Name", ""},
              {"ExpStr", json::parse("[\"undefined\"]")},
              {"ExpAct", "[]"_json},
              {"Params", "[]"_json},
              {"Args", "[]"_json},
              {"Condition", "[]"_json},
              {"ID", 41},
              {"First", "[]"_json},
              {"Second", "[]"_json},
              {"Degree", "[]"_json},
              {"Value", "[[]]"_json},
              {"Indexes", "[[]]"_json},
              {"Index_Type", ""},
              {"Hash_Values", "{}"_json},
              {"ValueType", json::parse("[" + type.dump() + "]")}
            };

            Returner fparsed = Returner{ noRet, vars, fRet, "none" };

            if (vars.find(name) == vars.end()) parsed = fparsed;
            else {

              json var = vars[name]["value"];

              for (json it : actions[i]["Indexes"]) {

                json _index = parser(it, calc_params, vars, dir, false, line, true).exp["ExpStr"][0];
                string index = _index.dump().substr(1, _index.dump().length() - 2);

                if (var["Hash_Values"].find(index) == var["Hash_Values"].end()) {
                  parsed = fparsed;
                  goto stopIndexing;
                }

                var = parser(var["Hash_Values"][index], calc_params, vars, dir, false, line, true).exp;
              }

              if (var["Type"] != "process") {
                parsed = fparsed;
                goto stopIndexing;
              }

              json params = var["Params"]
              , args = actions[i]["Args"];

              json sendVars = vars;

              for (int o = 0; o < params.size() || o < args.size(); o++) {

                json cur = {
                  {"type", "local"},
                  {"name", (string) params[o]},
                  {"value", parser(json::parse("[" + args[o].dump() + "]"), calc_params, vars, dir, false, line, true).exp},
                  {"valueActs", json::parse("[]")}
                };

                sendVars[(string) params[o]] = cur;
              }

              parsed = parser(var["ExpAct"], calc_params, sendVars, dir, true, line, false);

              json pVars = parsed.variables;

              //filter the variables that are not global
              for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                if (!(o.value()["type"] != "global" && o.value()["type"] != "process"))
                  vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();
            }

            stopIndexing:
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

              string val = (string) parser(actions[i]["Condition"][o]["Condition"], calc_params, vars, dir, false, line, true).exp[0][0];

              if (val != "false" && val != "undefined" && val != "null") {

                Returner parsed = parser(actions[i]["Condition"][o]["Actions"], calc_params, vars, dir, true, line, false);

                json pVars = parsed.variables;

                //filter the variables that are not global
                for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                  if (!(o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end()))
                    vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();

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

            string fileName = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0];

            if (fileName.rfind("\'", 0) == 0 || fileName.rfind("\"", 0) == 0 || fileName.rfind("`", 0) == 0) fileName = fileName.substr(1, fileName.length() - 2);

            string readerFile = dir + fileName
            , errMsg = "Could Not Find File: " + fileName;

            char* file = CReadFile(&readerFile[0], &errMsg[0], 1);

            string _acts = Cactions( CLex(file) );

            json acts = json::parse(_acts);

            Returner parsed = parser(acts, calc_params, vars, dir, false, 1, false);

            json pVars = parsed.variables;

            //filter the variables that are not global
            for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
              if (!(o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end()))
                vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();
          }
          break;
        case 15: {

            //read

            string in;

            cout << ((string) parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp["ExpStr"][0]) << " ";

            cin >> in;

            if (expReturn) {
              vector<string> retNo;

              json type = {
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
                {"Hash_Values", "{}"_json},
                {"ValueType", "[]"_json}
              }
              , expRet = {
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
                {"Hash_Values", "{}"_json},
                {"ValueType", json::parse("[" + type.dump() + "]")}
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

            string _code = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0].dump()
            , code = _code.substr(2, _code.length() - 4);

            char* codeNQ = NQReplaceC(&code[0]);

            char* len = CLex(codeNQ);

            char* __acts = Cactions(len);

            string _acts(__acts);

            json acts = json::parse(_acts);

            Returner parsed = parser(acts, calc_params, vars, dir, false, line, false);

            expStr[expStr.size() - 1].push_back(json::parse("[\"" + parsed.value[0] + "\"]")[0]);
          }
          break;
        case 19: {

            //typeof

            Returner parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true);

            json exp = parsed.exp;

            json type = exp["ValueType"];

            if (type.size() == 0) type = json::parse(
              R"(
                [{"Args":[],"Condition":[],"Degree":[],"ExpAct":[],"ExpStr":["type"],"First":[],"Hash_Values":{},"ID":44,"Index_Type":"","Indexes":[],"Name":"","Params":[],"Second":[],"Type":"type","Value":[],"ValueType":[]}]
              )"
            );

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars, type[0], "expression" };
          }
          break;
        case 20: {

            //err

            Returner parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true);

            string exp = parsed.exp[0][0].dump().substr(1, parsed.exp[0][0].dump().length() - 2);

            cout << exp << "\n\nerr~" << exp << "\n^^^^ <-- Error On Line " << line << endl;

            Kill();
          }
          break;
        case 21: {

            //loop

            json cond = actions[i]["Condition"][0]["Condition"]
            , acts = actions[i]["Condition"][0]["Actions"];

            Returner parsed;

            json condP = parser(cond, calc_params, vars, dir, false, line, true).exp[0][0];

            while (condP != "false" && condP != "undefined" && condP != "null") {

              parsed = parser(acts, calc_params, vars, dir, true, line, false);

              json pVars = parsed.variables;

              //filter the variables that are not global
              for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
                if (!(o.value()["type"] != "global" && o.value()["type"] != "process" && vars.find(o.value()["name"]) != vars.end()))
                  vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") continue;
              if (parsed.type == "break") break;

              condP = parser(cond, calc_params, vars, dir, false, line, true).exp[0][0];
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

            cout << index.dump(2) << endl;

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

                json ascValType = {
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
                  {"Hash_Values", "{}"_json},
                  {"ValueType", "[]"_json}
                }
                , ascVal = {
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
                  {"Hash_Values", "{}"_json},
                  {"ValueType", json::parse("[" + ascValType.dump() + "]")}
                };

                return Returner{returnNone, vars, ascVal, "expression"};
              }
            }
          }
          break;
        case 27: {

            //parse

            json parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            vector<string> returnNone;

            if ((parsed["Type"] != "string" && parsed["Type"] != "number") && expReturn) return Returner{ returnNone, vars, falseyVal, "expression" };
            else {

              string putVal;

              if (parsed["Type"] == "number") putVal = "\"" + parsed["ExpStr"][0].get<string>() + "\"";
              else if (parsed["Type"] == "string") putVal = "\"" + parsed["ExpStr"][0].get<string>().substr(1, parsed["ExpStr"][0].get<string>().length() - 2) + "\"";

              json typeVal = {
                {"Type", "number"},
                {"Name", ""},
                {"ExpStr", json::parse("[" + putVal + "]")},
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
                {"Hash_Values", "{}"_json},
                {"ValueType", "[]"_json}
              }
              , retVal = {
                {"Type", "number"},
                {"Name", ""},
                {"ExpStr", json::parse("[" + putVal + "]")},
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
                {"Hash_Values", "{}"_json},
                {"ValueType", json::parse("[" + typeVal.dump() + "]")}
              };

              return Returner{ returnNone, vars, retVal, "expression" };
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
        case 31: {

            //len

            json calculated = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp;

            json lenPlaceholder = {
              {"Type", "number"},
              {"Name", ""},
              {"ExpStr", json::parse("[\"0\"]")},
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
              {"Hash_Values", "{}"_json},
              {"ValueType", {
                {
                  {"Type", "number"},
                  {"Name", ""},
                  {"ExpStr", json::parse("[\"0\"]")},
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
                  {"Hash_Values", "{}"_json},
                  {"ValueType", "[]"_json}
                }
              }}
            };

            if (expReturn) {

              vector<string> returnNone;

              if (calculated.size() == 0) return Returner{ returnNone, vars, lenPlaceholder, "expression" };

              switch (GetActNumC(&(calculated["Type"].get<string>())[0])) {
                case 38: {

                    json nExpStr = json::parse("[\"" + to_string(calculated["ExpStr"][0].get<string>().length() - 2) + "\"]");

                    lenPlaceholder["ExpStr"] = nExpStr;
                    lenPlaceholder["ValueType"][0]["ExpStr"] = nExpStr;
                  }
                  break;
                case 39: {

                    json nExpStr = calculated["ExpStr"];

                    lenPlaceholder["ExpStr"] = nExpStr;
                    lenPlaceholder["ValueType"][0]["ExpStr"] = nExpStr;
                  }
                  break;
                case 40:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 41:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 42:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 44: {

                    json nExpStr = json::parse("[\"" + to_string(calculated["ExpStr"][0].get<string>().length() - 2) + "\"]");

                    lenPlaceholder["ExpStr"] = nExpStr;
                    lenPlaceholder["ValueType"][0]["ExpStr"] = nExpStr;
                  }
                case 45:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 46:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 10:
                  return Returner{ returnNone, vars, lenPlaceholder, "expression" };
                  break;
                case 22: {

                    json nExpStr = json::parse("[\"" + to_string(calculated["Hash_Values"].size()) + "\"]");

                    lenPlaceholder["ExpStr"] = nExpStr;
                    lenPlaceholder["ValueType"][0]["ExpStr"] = nExpStr;
                  }
                  break;
                case 24: {

                    json nExpStr = json::parse("[\"" + to_string(calculated["Hash_Values"].size()) + "\"]");

                    lenPlaceholder["ExpStr"] = nExpStr;
                    lenPlaceholder["ValueType"][0]["ExpStr"] = nExpStr;
                  }
                  break;
                default: return Returner{ returnNone, vars, lenPlaceholder, "expression" };
              }

              return Returner{ returnNone, vars, lenPlaceholder, "expression" };
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

            json var = vars[actions[i]["Name"].get<string>()]["value"];

            bool varIsMutable = var["IsMutable"].get<bool>()
            , actIsMutable = actions[i]["IsMutable"].get<bool>()
            , isMutable = varIsMutable ^ actIsMutable;

            var["IsMutable"] = isMutable;

            json val = parser(json::parse("[" + var.dump() + "]"), calc_params, vars, dir, false, line, true).exp;

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

        //assignment operators
        case 4343: {

            //++

            string name = actions[i]["Name"];

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Add(&(val)[0], "1", &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"1\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 4545: {

            //--

            string name = actions[i]["Name"];

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Subtract(&(val)[0], "1", &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"1\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 4361: {

            //+=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Add(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 4561: {

            //-=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Subtract(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 4261: {

            //*=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Multiply(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 4761: {

            ///=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Division(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 9461: {

            //^=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Exponentiate(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
          }
          break;
        case 3761: {

            //%=

            string name = actions[i]["Name"];

            json __inc = actions[i]["ExpAct"]
            , _inc = parser(__inc, calc_params, vars, dir, false, line, true).exp[0][0];
            string inc = _inc.dump().substr(1, _inc.dump().length() - 2);

            json nVar;

            if (vars[name]["type"] != "dynamic") {

              if (vars[name].find("value") != vars[name].end()) {

                string _val = vars[name]["value"][0].dump()
                , val = _val.substr(2, _val.length() - 4);

                char* _added = Modulo(&(val)[0], &inc[0], &calc_params.dump()[0], line);
                string added(_added);

                nVar = {
                  {"type", vars[name]["type"]},
                  {"name", name},
                  {"value", json::parse("[[\"" + added + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
              } else
                nVar = {
                  {"type", "local"},
                  {"name", name},
                  {"value", json::parse("[[\"" + inc + "\"]]")},
                  {"valueActs", json::parse("[]")}
                };
            }

            vars[name] = nVar;
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
  ret.exp = "{}"_json;
  ret.type = "none";

  return ret;
}

#endif
