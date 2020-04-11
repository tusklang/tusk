#ifndef PARSER_HPP_
#define PARSER_HPP_

#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include "json.hpp"
#include "bind.h"
#include "math.hpp"
#include "structs.h"
#include "indexes.hpp"
#include "log_format.hpp"
using namespace std;
using json = nlohmann::json;

json math(json exp, const json calc_params, json vars, const string dir, int line);

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

            json index = indexesCalc(actions[i], actions[i]["Indexes"], calc_params, vars, line, dir, actions[i]["Index_Type"]);

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

            Returner parsed = parser(var["valueActs"], calc_params, sendVars, dir, true, line, false);

            json pVars = parsed.variables;

            //filter the variables that are not global
            for (json::iterator o = pVars.begin(); o != pVars.end(); o++)
              if (!(o.value()["type"] != "global" && o.value()["type"] != "process"))
                vars[o.value()["name"].dump().substr(1, o.value()["name"].dump().length() - 2)] = o.value();

            expStr[expStr.size() - 1].push_back(parsed.value[0]);
          }
          break;
        case 12:

          //return

          return Returner{ parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0], vars, expStr, "return" };
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

            cout << ((string) parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0]) << " ";

            cin >> in;

            expStr[expStr.size() - 1].push_back(json::parse("[\"\'" + in + "\'\"]")[0]);
          }
          break;
        case 16: {

            //break

            Returner ret;

            vector<string> returnNone;

            ret.value = returnNone;
            ret.variables = vars;
            ret.exp = math(expStr, calc_params, vars, dir, line);
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
            ret.exp = math(expStr, calc_params, vars, dir, line);
            ret.type = "skip";

            return ret;
          }
          break;
        case 18: {

            //eval

            string _code = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0].dump()
            , code = _code.substr(2, _code.length() - 4);

            char* codeNQ = NQReplaceC(&code[0]);

            char* lex = CLex(codeNQ);

            char* __acts = Cactions(lex);

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

              return Returner{ returnNone, vars, actions[i], "expression" };
            }
          }
          break;
        case 23: {

            //hashIndex

            json val = actions[i]["Hash_Values"];

            json index = indexesCalc(val, actions[i]["Indexes"], calc_params, vars, line, dir, actions[i]["Index_Type"]);

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

              return Returner{ returnNone, vars, actions[i], "expression"};
            }
          }
          break;
        case 25: {

            //arrayIndex

            json val = actions[i]["Value"]
            , index = indexesCalc(val, actions[i]["Indexes"], calc_params, vars, line, dir, actions[i]["Index_Type"]);

            if (expReturn) {
              vector<string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 26: {

            //ascii

            string parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0].dump();

            parsed = parsed.substr(1, parsed.length() - 2);

            char first = parsed[0];

            expStr[expStr.size() - 1].push_back(json::parse("[\"" + to_string((int) first) + "\"]")[0]);
          }
          break;
        case 27: {

            //parse

            string parsed = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0][0].dump();

            parsed = parsed.substr(1, parsed.length() - 2);

            if (!(strcmp(GetType(&parsed[0]), "string") == 0 || strcmp(GetType(&parsed[0]), "number") == 0)) {
              cout << "There Was An Error: `parse~` cannot be used on a non-string or number" << "\n\nparse~" + parsed << "\n^ <-- Error On Line " + line;
              Kill();
            }

            if (strcmp(GetType(&parsed[0]), "string") == 0)
              expStr[expStr.size() - 1].push_back(json::parse("[\"" + parsed.substr(1, parsed.length() - 2) + "\"]")[0]);
            else expStr[expStr.size() - 1].push_back(json::parse("[\"" + parsed + "\"]")[0]);
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

            if (vars[name].find("type") != vars[name].end())
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

            vars[name] = nVar;
          }
          break;
        case 31: {

            //len

            json calculated = parser(actions[i]["ExpAct"], calc_params, vars, dir, false, line, true).exp[0];

            json noNewline;

            copy_if(calculated.begin(), calculated.end(), back_inserter(noNewline), [](json i){
              return i != "newlineN";
            });

            string datatype_get = "";

            for (string o : noNewline) datatype_get+=o;

            string type(GetType(&datatype_get[0]));

            //TODO: maybe switch to a switch statement later

            if (type == "string") expStr[expStr.size() - 1].push_back( to_string(((string) datatype_get).length() - 2) );
            else if (type == "hash") {
              int commas = 0;

              int bCnt = 0;

              for (int o = 0; o < noNewline.size(); o++) {

                json it = noNewline[o];

                if (it == "[:" || it == "[") bCnt++;
                if (it == ":]" || it == "]") bCnt--;

                if (bCnt == 1 && it == ",") commas++;
              }

              expStr[expStr.size() - 1].push_back( to_string(commas + 1) );
            } else if (type == "array") {
              int commas = 0;

              int bCnt = 0;

              for (int o = 0; o < noNewline.size(); o++) {

                json it = noNewline[o];

                if (it == "[:" || it == "[") bCnt++;
                if (it == ":]" || it == "]") bCnt--;

                if (bCnt == 1 && it == ",") commas++;
              }

              expStr[expStr.size() - 1].push_back( to_string(commas + 1) );
            } else if (type == "boolean") expStr[expStr.size() - 1].push_back(datatype_get == "true" ? "1" : "0");
            else if (type == "falsey") expStr[expStr.size() - 1].push_back("0");
            else if (type == "number") expStr[expStr.size() - 1].push_back(datatype_get);
            else expStr[expStr.size() - 1].push_back("0");

          }
          break;
        case 32: {

            //add

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp.dump()
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp.dump();

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

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp[0][0]
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp[0][0];

            expStr[expStr.size() - 1].push_back(
              Subtract(
                &first[0],
                &second[0],
                &calc_params.dump()[0],
                line
              )
            );
          }
          break;
        case 34: {

            //multiply

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp[0]
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp[0];

            expStr[expStr.size() - 1].push_back(
              Multiply(
                &first[0],
                &second[0],
                &calc_params.dump()[0],
                line
              )
            );
          }
          break;
        case 35: {

            //divide

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp[0][0]
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp[0][0];

            expStr[expStr.size() - 1].push_back(
              Division(
                &first[0],
                &second[0],
                &calc_params.dump()[0],
                line
              )
            );
          }
          break;
        case 36: {

            //exponentiate

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp[0][0]
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp[0][0];

            expStr[expStr.size() - 1].push_back(
              Exponentiate(
                &first[0],
                &second[0],
                &calc_params.dump()[0],
                line
              )
            );
          }
          break;
        case 37: {

            //modulo

            string first = parser(actions[i]["First"], calc_params, vars, dir, false, line, true).exp[0][0]
            , second = parser(actions[i]["Second"], calc_params, vars, dir, false, line, true).exp[0][0];

            expStr[expStr.size() - 1].push_back(
              Modulo(
                &first[0],
                &second[0],
                &calc_params.dump()[0],
                line
              )
            );
          }
          break;
        case 38: {

            //string

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars,  actions[i], "expression" };
          }
          break;
        case 39: {

            //number

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars,  actions[i], "expression" };
          }
          break;
        case 40: {

            //boolean

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars,  actions[i], "expression" };
          }
          break;
        case 41: {

            //falsey

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars,  actions[i], "expression" };
          }
          break;
        case 42: {

            //none

            vector<string> noRet;

            if (expReturn) return Returner{ noRet, vars,  actions[i], "expression" };
          }
          break;
        case 43: {

            //variable

            if (expReturn) return parser(actions[i], calc_params, vars, dir, false, line, true);
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
  ret.exp = math(expStr, calc_params, vars, dir, line);
  ret.type = "none";

  return ret;
}

#endif
