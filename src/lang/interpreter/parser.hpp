#ifndef PARSER_HPP_
#define PARSER_HPP_

#include <iostream>
#include <deque>
#include <vector>
#include <map>
#include <string>
#include <algorithm>
#include <thread>
#include <windows.h>
#include <regex>
#include <exception>
#include "json.hpp"
#include "../bind.h"
#include "structs.hpp"
#include "indexes.hpp"
#include "log_format.hpp"
#include "values.hpp"
#include "comparisons.hpp"
#include "similarity.hpp"
#include "processes.hpp"
#include "ommtypes.hpp"

//file i/o
#include "../files/readfile.hpp"
#include "../files/writefile.h"
#include "../files/delete.hpp"
#include "../files/isDir.hpp"
#include "../files/isFile.hpp"
#include "../files/readdir.hpp"

//operations
#include "operations/add/add.hpp"
#include "operations/divide/divide.hpp"
#include "operations/exponentiate/exponentiate.hpp"
#include "operations/modulo/modulo.hpp"
#include "operations/multiply/multiply.hpp"
#include "operations/subtract/subtract.hpp"

using json = nlohmann::json;

namespace omm {
  Returner parser(const std::vector<Action> actions, const json cli_params, std::map<std::string, Variable> vars, const bool groupReturn, const bool expReturn, deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    //loop through every action
    for (Action v : actions) {

      //get current action id
      int cur = v.ID;

      switch (cur) {
        case 1: {

            //local

            std::string name = v.Name;

            std::vector<Action> acts = v.ExpAct
            , parsed = { parser(acts, cli_params, vars, false, true, this_vals, dir).exp };

            Variable nVar = Variable{
              "local",
              name,
              parsed
            };

            vars[name] = nVar;
          }
          break;
        case 2: {

            //dynamic

            std::string name = v.Name;

            std::vector<Action> acts = v.ExpAct;

            Variable nVar = Variable{
              "dynamic",
              name,
              acts
            };

            vars[name] = nVar;
          }
          break;
        case 3: {

            //alt

            int o = 0;

            Returner cond = parser(v.Condition[0].Condition, cli_params, vars, false, true, this_vals, dir);

            //while the alt statement should continue
            while (isTruthy(cond.exp)) {

              //going back to the first block when it reached the last block
              if (o >= v.Condition.size()) o = 0;

              Returner parsed = parser(v.Condition[o].Actions, cli_params, vars, true, false, this_vals, dir);

              std::map<std::string, Variable> pVars = parsed.variables;

              //filter the variables that are not global
              for (std::pair<std::string, Variable> o : pVars)
                if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                  vars[o.first] = o.second;

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") continue;
              if (parsed.type == "break") break;

              cond = parser(v.Condition[o].Condition, cli_params, vars, false, true, this_vals, dir);
              o++;
            }
          }
          break;
        case 4: {

            //global

            std::string name = v.Name;

            std::vector<Action> acts = v.ExpAct
            , parsed = { parser(acts, cli_params, vars, false, true, this_vals, dir).exp };

            Variable nVar = Variable{
              "global",
              name,
              parsed
            };

            vars[name] = nVar;
          }
          break;
        case 5: {

            //log

            Action _val = parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir).exp;

            log_format(_val, cli_params, vars, 2, "log");
          }
          break;
        case 6: {

            //print

            Action _val = parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir).exp;

            log_format(_val, cli_params, vars, 2, "print");
          }
          break;
        case 8: {

            //expressionIndex

            Action val = parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir).exp;

            Action index = indexesCalc(val.Hash_Values, v.Indexes, cli_params, vars, this_vals, dir);

            if (expReturn) {
              std::vector<std::string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 9: {

            //group

            std::vector<Action> acts = v.ExpAct;

            Returner parsed = parser(acts, cli_params, vars, false, false, this_vals, dir);

            std::map<std::string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (std::pair<std::string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                vars[o.first] = o.second;

            if (groupReturn || expReturn) return Returner{ parsed.value, vars, parsed.exp, parsed.type };
          }
          break;
        case 10: {

          //process

          std::string name = v.Name;

          if (name != "") {
            Variable nVar = Variable{
              "process",
              name,
              { v }
            };

            vars[name] = nVar;
          }

          if (expReturn) {
            std::vector<std::string> noRet;

            return Returner{ noRet, vars, v, "expression" };
          }

          break;
        }
        case 80: {

          //pargc_number

          unsigned long long pargc = 0;

          //count the pargc
          for (std::pair<std::string, Variable> it : vars)
            if (it.second.type == "argument") ++pargc;
            else if (it.second.type == "pargv")
              pargc+=parser(it.second.value, cli_params, vars, false, true, this_vals, dir).exp.Hash_Values.size();

          if (std::to_string(pargc) == std::string(ReturnInitC(&v.ExpStr[0][0]))) {

            Returner parsed = parser(v.ExpAct, cli_params, vars, true, true, this_vals, dir);

            std::map<std::string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (std::pair<std::string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                vars[o.first] = o.second;

            if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
            if (parsed.type == "skip") continue;
            if (parsed.type == "break") break;
          }

          break;
        }
        case 81: {

          //pargc_paramlist

          unsigned long long pargc = 0;
          std::vector<std::string> types;

          //count the pargc and the types
          for (std::pair<std::string, Variable> it : vars) {

            Action parsed_it = parser(it.second.value, cli_params, vars, false, true, this_vals, dir).exp;

            if (it.second.type == "argument") {

              ++pargc;
              types.push_back(parsed_it.Type);

            } else if (it.second.type == "pargv") {

              pargc+=parsed_it.Hash_Values.size();

              for (std::pair<std::string, std::vector<Action>> pargv_it: parsed_it.Hash_Values)
                types.push_back(parser(pargv_it.second, cli_params, vars, false, true, this_vals, dir).exp.Type);
            }
          }

          if (pargc == v.Params.size() && types == v.Params) {
            Returner parsed = parser(v.ExpAct, cli_params, vars, true, true, this_vals, dir);

            std::map<std::string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (std::pair<std::string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                vars[o.first] = o.second;

            if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
            if (parsed.type == "skip") continue;
            if (parsed.type == "break") break;
          }

          break;
        }
        case 11: {

            //# (call process)

            std::string name = v.Name;

            Returner parsed;

            std::vector<std::string> noRet;

            Returner fparsed = Returner{ noRet, vars, falseyVal, "expression" };

            parsed = fparsed;

            if (vars.find(name) == vars.end()) goto stopIndexing_processes;
            else {

              Action var = vars[name].value[0];

              parsed = processParser(var, v, cli_params, &vars, this_vals, true, dir);
            }

            stopIndexing_processes:
            if (expReturn) {

              Action val = parsed.exp;

              std::vector<std::string> noRet;

              return Returner{ noRet, vars, val, "expression" };
            }
          }
          break;
        case 12: {

            //return

            std::vector<std::string> noRet;

            return Returner{ noRet, vars, parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir).exp, "return" };
          }
          break;
        case 13: {

            //conditional

            for (int o = 0; o < v.Condition.size(); o++) {

              Action val = parser(v.Condition[o].Condition, cli_params, vars, false, true, this_vals, dir).exp;

              if (isTruthy(val) || v.Condition[o].Type == "else") {

                Returner parsed = parser(v.Condition[o].Actions, cli_params, vars, true, false, this_vals, dir);

                std::map<std::string, Variable> pVars = parsed.variables;

                //filter the variables that are not global
                for (std::pair<std::string, Variable> o : pVars)
                  if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                    vars[o.first] = o.second;

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

            std::vector<std::vector<Action>> files = v.Value; //get all actionized files imported

            //loop through actionized files
            for (std::vector<Action> it : files) {

              Returner parsed = parser(it, cli_params, vars, true, false, this_vals, dir);

              std::map<std::string, Variable> pVars = parsed.variables;

              //filter the variables that are not global
              for (std::pair<std::string, Variable> o : pVars)
                if (o.second.type == "global" || o.second.type == "process")
                  vars[o.first] = o.second;
            }
          }
          break;
        case 16: {

            //break

            Returner ret;

            std::vector<std::string> returnNone;
            Action expNone;

            ret.value = returnNone;
            ret.variables = vars;
            ret.exp = expNone;
            ret.type = "break";

            return ret;
          }
          break;
        case 17: {

            //skip

            Returner ret;

            std::vector<std::string> returnNone;
            Action expNone;

            ret.value = returnNone;
            ret.variables = vars;
            ret.exp = expNone;
            ret.type = "skip";

            return ret;
          }
          break;
        case 21: {

            //loop

            std::vector<Action> cond = v.Condition[0].Condition
            , acts = v.Condition[0].Actions;

            Returner parsed;

            Action condP = parser(cond, cli_params, vars, false, true, this_vals, dir).exp;

            while (isTruthy(condP)) {

              parsed = parser(acts, cli_params, vars, true, false, this_vals, dir);

              std::map<std::string, Variable> pVars = parsed.variables;

              //filter the variables that are not global
              for (std::pair<std::string, Variable> o : pVars)
                if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                  vars[o.first] = o.second;

              if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
              if (parsed.type == "skip") continue;
              if (parsed.type == "break") break;

              condP = parser(cond, cli_params, vars, false, true, this_vals, dir).exp;
            }

          }
          break;
        case 22: {

            //hash

            if (expReturn) {

              std::vector<std::string> returnNone;

              bool isMutable = v.IsMutable;

              Action val = v;

              if (!isMutable) {

                for (std::pair<std::string, std::vector<Action>> it : v.Hash_Values) {

                  std::string paramCount = "";
                  Action exp = parser(it.second, cli_params, vars, false, true, this_vals, dir).exp;

                  val.Hash_Values[it.first] = { exp };
                }
              }

              return Returner{ returnNone, vars, val, "expression" };
            }
          }
          break;
        case 23: {

            //hashIndex

            std::map<std::string, std::vector<Action>> val = v.Hash_Values;

            Action index = indexesCalc(val, v.Indexes, cli_params, vars, this_vals, dir);

            if (expReturn) {
              std::vector<std::string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 24: {

            //array

            if (expReturn) {
              std::vector<std::string> returnNone;

              bool isMutable = v.IsMutable;

              Action val = v;

              if (!isMutable) {

                char* index = "0";

                for (std::pair<std::string, std::vector<Action>> o : v.Hash_Values) {

                  if (val.Hash_Values.find(index) == val.Hash_Values.end()) {
                    index = AddC(index, "1", &cli_params.dump()[0]);
                    continue;
                  }

                  std::string paramCount = "";
                  Action exp = parser(o.second, cli_params, vars, false, true, this_vals, dir).exp;

                  val.Hash_Values[index] = { exp };
                  index = AddC(index, "1", &cli_params.dump()[0]);
                }
              }

              return Returner{ returnNone, vars, val, "expression"};
            }
          }
          break;
        case 25: {

            //arrayIndex

            std::map<std::string, std::vector<Action>> val = v.Hash_Values;
            Action index = indexesCalc(val, v.Indexes, cli_params, vars, this_vals, dir);

            if (expReturn) {
              std::vector<std::string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 28: {

            //let

            std::string name = v.Name;

            std::vector<Action> acts = v.ExpAct;

            std::vector<Action> parsed = { parser(acts, cli_params, vars, false, true, this_vals, dir).exp };

            Variable nVar;

            Variable var = vars[name];
            std::vector<std::string> indexes;

            if (v.Indexes.size() == 0) {

              if (vars.find(name) != vars.end())
                vars[name] = Variable{
                  vars[name].type,
                  name,
                  parsed
                };
              else
                vars[name] = Variable{
                  "local",
                  name,
                  parsed
                };
            } else {

              Action* map = &vars[name].value[0];

              for (std::vector<Action> it : v.Indexes) {

                std::string varP = parser(it, cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

                map = &(map->Hash_Values[varP][0]);
              }

              *map = parsed[0];

            }
          }
          break;
        case 32: {

            //add

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = add(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = subtract(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = multiply(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = divide(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = exponentiate(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = modulo(first, second, cli_params, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = val;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 38: {

            //std::string

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, v, "expression" };
          }
          break;
        case 39: {

            //number

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, v, "expression" };
          }
          break;
        case 40: {

            //boolean

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, v, "expression" };
          }
          break;
        case 41: {

            //falsey

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, v, "expression" };
          }
          break;
        case 42: {

            //none

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, v, "expression" };
          }
          break;
        case 43: {

            //variable

            Action val;

            if (vars.find(v.Name) == vars.end()) val = falseyVal;
            else {

              Action var = vars[v.Name].value[0];

              bool varIsMutable = var.IsMutable
              , actIsMutable = v.IsMutable
              , isMutable = varIsMutable ^ actIsMutable;

              var.IsMutable = isMutable;

              val = parser({ var }, cli_params, vars, false, true, this_vals, dir).exp;
            }

            std::vector<std::string> noRet;

            if (expReturn) return Returner{ noRet, vars, val, "expression" };
          }
          break;
        case 46: {

            //variableIndex

            Returner parsedVal = parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir);

            Action index = indexesCalc(parsedVal.exp.Hash_Values, v.Indexes, cli_params, vars, this_vals, dir);

            if (expReturn) {
              std::vector<std::string> returnNone;

              return Returner{ returnNone, vars, index, "expression" };
            }
          }
          break;
        case 47: {

            //equals

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action val = equals(
              first,
              second,
              cli_params,
              vars,
              this_vals,
              dir
            );

            if (first.Type != second.Type) val = falseRet;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          Action val = equals(
            first,
            second,
            cli_params,
            vars,
            this_vals,
            dir
          );

          val = val.ExpStr[0] == "true" ? falseRet : trueRet;
          if (first.Type != second.Type) val = trueRet;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 49: {

          //greater

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          Action val = isGreater(
            first,
            second,
            cli_params
          );

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 50: {

          //less

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          Action val = isLess(
            first,
            second,
            cli_params
          );

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 51: {

          //greaterOrEqual

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          Action val = isLess(
            first,
            second,
            cli_params
          );

          val = val.ExpStr[0] == "true" ? falseRet : trueRet;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 52: {

          //lessOrEqual

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          Action val = isGreater(
            first,
            second,
            cli_params
          );

          val = val.ExpStr[0] == "true" ? falseRet : trueRet;
          if (first.Type != second.Type) val = trueRet;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = val;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 53: {

            //not

            Action val = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp
            , retval;

            std::string expstr = val.ExpStr[0];

            if (expstr == "false" || val.Type == "falsey") retval = trueRet;
            else retval = falseRet;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action retval;

            if ((v.Degree).size() == 0) retval = similarity(first, second, zero, cli_params, vars, this_vals, dir);
            else retval = similarity(first, second, parser(v.Degree, cli_params, vars, false, true, this_vals, dir).exp, cli_params, vars, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
            , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

            Action retval;

            if ((v.Degree).size() == 0) retval = strictSimilarity(first, second, zero, cli_params, vars, this_vals, dir);
            else retval = strictSimilarity(first, second, parser(v.Degree, cli_params, vars, false, true, this_vals, dir).exp, cli_params, vars, this_vals, dir);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = retval;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 71: {

          //or

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = isTruthy(first) || isTruthy(second) ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 72: {

          //and

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = isTruthy(first) && isTruthy(second) ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 73: {

          //nor

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = isTruthy(first) || isTruthy(second) ? falseRet : trueRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 74: {

          //nand

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = isTruthy(first) && isTruthy(second) ? falseRet : trueRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 75: {

          //xor

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = (isTruthy(first) || isTruthy(second)) && !(isTruthy(first) && isTruthy(second)) ? trueRet : falseRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 76: {

          //xnor

          Action first = parser(v.First, cli_params, vars, false, true, this_vals, dir).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals, dir).exp;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = (isTruthy(first) || isTruthy(second)) && !(isTruthy(first) && isTruthy(second)) ? falseRet : trueRet;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 56:  {

            //@ (call thread)

            std::string name = v.Name;

            Returner parsed;

            std::vector<std::string> noRet;

            Returner fparsed = Returner{ noRet, vars, falseyVal, "none" };

            parsed = fparsed;

            if (vars.find(name) == vars.end()) goto stopIndexing_threads;
            else {

              Action var = vars[name].value[0];

              for (std::vector<Action> it : v.Indexes) {

                std::string index = parser(it, cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

                if (var.Hash_Values.find(index) == var.Hash_Values.end()) {
                  parsed = fparsed;
                  goto stopIndexing_threads;
                }

                var = parser(var.Hash_Values[index], cli_params, vars, false, true, this_vals, dir).exp;
              }

              if (var.Type != "process") {
                parsed = fparsed;
                goto stopIndexing_threads;
              }

              std::vector<std::string> params = var.Params;
              std::vector<std::vector<Action>> args = v.Args;

              if (params.size() != args.size() && !vector_indexes_inc(params, "pargv")) {
                parsed = fparsed;
                return Returner{ noRet, vars, falseyVal, "expression" };
              }

              std::map<std::string, Variable> sendVars = vars;

              for (int o = 0; o < params.size() || o < args.size(); o++) {

                //if it starts with pargv
                if (params[o].rfind("$pargv.", 0) == 0) {

                  std::string varname = "$" + params[o].substr(std::string("$pargv.").length());

                  //convert the rest of the args into an array and store it in the pargv variable
                  std::map<std::string, std::vector<Action>> pargv;

                  for (unsigned long long cur = 0; o < args.size(); ++o, ++cur)
                    pargv[std::to_string(cur)] = { parser(args[o], cli_params, vars, false, true, this_vals, dir).exp };

                  Action arg = arrayVal;
                  arg.Hash_Values = pargv;

                  sendVars[varname] = Variable{
                    "pargv",
                    varname,
                    { arg }
                  };

                  break;
                }

                Variable cur = Variable{
                  "argument",
                  params[o],
                  { parser(args[o], cli_params, vars, false, true, this_vals, dir).exp }
                };

                sendVars[params[o]] = cur;
              }

              std::thread _(parser, var.ExpAct, cli_params, sendVars, true, false, this_vals, dir);

              _.detach();
            }

            stopIndexing_threads:
            if (expReturn) {

              Action val = parsed.exp;

              std::vector<std::string> noRet;

              return Returner{ noRet, vars, val, "expression" };
            }
          }
          break;
        case 57: {

            //wait

            Action amt = parser(v.ExpAct, &cli_params.dump()[0], vars, false, true, this_vals, dir).exp;

            if (IsLessC(&(amt.ExpStr[0])[0], "4294967296")) Sleep((unsigned long) std::atoi(&(amt.ExpStr[0])[0]));
            else {
              for (char* o = "0"; (bool) IsLessC(o, &(amt.ExpStr[0])[0]); o = AddC(o, "4294967296", &cli_params.dump()[0])) {

                char* subtracted = SubtractC(&(amt.ExpStr[0])[0], o, &cli_params.dump()[0]);

                if (IsLessC(
                  subtracted,
                  "4294967296"
                )) Sleep((unsigned long) std::atoi(subtracted));
                else Sleep((unsigned long) 4294967296);
              }
            }
          }
          break;
        case 58: {

            //cast

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

              Action cur = parser(v.ExpAct, cli_params, vars, false, true, this_vals, dir).exp;
              cur.Type = v.Name;
              cur.Name = v.ExpStr[0];

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = cur;
              ret.type = "expression";

              return ret;
            }
          }
          break;
        case 59: {

            //each

            std::vector<std::string> putterVars = v.ExpStr;
            std::string var1 = putterVars[0]
            , var2 = putterVars[1];

            //parse the iterator value
            std::map<std::string, std::vector<Action>> iterator = parser(v.First /* v.First is where the iterator is stored */, cli_params, vars, false, true, this_vals, dir).exp.Hash_Values;

            iterator.erase("falsey");

            for (std::pair<std::string, std::vector<Action>> it : iterator) {
              std::map<std::string, Variable> sendVars = vars;

              Action key = strPlaceholder;

              sendVars[var1] = Variable{
                "local",
                var1,
                { key }
              };
              sendVars[var2] = Variable{
                "local",
                var2,
                { parser(it.second, cli_params, vars, false, true, this_vals, dir).exp }
              };

              Returner parsed = parser(v.ExpAct, cli_params, sendVars, true, false, this_vals, dir);

              std::map<std::string, Variable> pVars = parsed.variables;

              //filter the variables that are not global
              for (std::pair<std::string, Variable> o : pVars)
                if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                  vars[o.first] = o.second;

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

          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          if (isDir(nDir + filename)) {

            std::vector<std::string> dirs = read_dir(nDir + filename);

            Action dir_arr = arrayVal;

            unsigned long long i = 0;
            for (std::string curD : dirs) {

              Action curDirOmm = ommtypes::to_string(curD);
              dir_arr.Hash_Values[std::to_string(i)] = { curDirOmm };

              ++i;
            }

            if (expReturn) {

              Returner ret;
              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = dir_arr;
              ret.type = "expression";

              return ret;
            }

          } else {
            std::string content = readfile(&(nDir + filename)[0]);

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

              Action contentJ = strPlaceholder;

              contentJ.ExpStr = {content};

              //make the hash values of the std::string
              for (unsigned long long o = 0; o < content.length(); o++) {
                Action curChar = strPlaceholder;

                curChar.ExpStr = {
                  std::string(1, content[o])
                };

                contentJ.Hash_Values[std::to_string(o)] = { curChar };
              }

              Action returner = contentJ;

              if (v.SubCall.size() > 0) {

                Action callProcessParser = v;

                bool isProc = v.SubCall[0].IsProc;

                callProcessParser.Indexes = v.SubCall[0].Indexes;
                callProcessParser.Args = v.SubCall[0].Args;
                callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

                returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

              }

              if (expReturn) {

                Returner ret;
                std::vector<std::string> retNo;

                ret.value = retNo;
                ret.variables = vars;
                ret.exp = returner;
                ret.type = "expression";

                return ret;
              }
            }
          }

          break;
        }
        case 61: {

          //files.write

          //written as files.write(dir, content)

          //get both arguments and parse them
          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];
          Action content = parser(v.Args[1], cli_params, vars, false, true, this_vals, dir).exp;

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          std::string contentstr = content.ExpStr[0];
          writefile(&(nDir + filename)[0], &contentstr[0]);

          Action returner = content;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

          }

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returner;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 78: {

          //files.remove

          //written as files.remove(dir)

          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          deletefile(nDir + filename);

          break;
        }
        case 62: {

          //files.exists

          //written as files.exists(dir)

          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          //if it is not a directory and not a file, it does not exist
          bool exists = !(!isDir(nDir + filename) && !isFile(nDir + filename));

          Action returner = exists ? trueRet : falseRet;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

          }

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returner;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 63: {

          //files.isFile

          //written as files.isFile(dir)

          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          bool isFileVal = isFile(nDir + filename);

          Action returner = isFileVal ? trueRet : falseRet;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

          }

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returner;
            ret.type = "expression";

            return ret;
          }
          break;
        }
        case 64: {

          //files.isDir

          //written as files.isDir(dir)

          std::string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          smatch match;

          //see if the filename is absolute
          std::regex pat("^[a-zA-Z]:");
          bool isOnDrive = std::regex_search(filename, match, pat);

          std::string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

          bool isDirVal = isDir(nDir + filename);

          Action returner = isDirVal ? trueRet : falseRet;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

          }

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returner;
            ret.type = "expression";

            return ret;
          }
          break;
        }

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

        case 68: {

          //regex.match

          std::string str = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];
          std::string regstr = parser(v.Args[1], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          try {
            std::regex reg(regstr);

            smatch matcher;

            std::vector<unsigned long long> found_indexes;

            //get all matches
            for (auto it = std::sregex_iterator(str.begin(), str.end(), reg); it != std::sregex_iterator(); it++) {
              found_indexes.push_back(it->position());
            }

            Action returnerArr = arrayVal;

            char* cur = "0";

            //loop through the indexes found and store them an omm type array
            for (int i : found_indexes) {

              //store the value of the number 1
              Action indexJ = val1;

              indexJ.ExpStr[0] = to_string(i);

              returnerArr.Hash_Values[std::string(cur)] = { indexJ };
              cur = AddC(cur, "1", &cli_params.dump()[0]);
            }

            Action returnerVal = hashVal;

            if (v.SubCall.size() > 0) {

              Action callProcessParser = v;

              bool isProc = v.SubCall[0].IsProc;

              callProcessParser.Indexes = v.SubCall[0].Indexes;
              callProcessParser.Args = v.SubCall[0].Args;
              callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

              returnerVal = processParser(returnerArr, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

            }

            if (expReturn) {

              Returner ret;
              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = returnerVal;
              ret.type = "expression";

              return ret;
            }

          } catch (std::regex_error& e) {

            //give information about the warning
            cout << "Warning during interpreting: Invalid Regular Expression: " << regstr << endl;
            cout << "Error description: " << e.what() << endl;
            cout << "Error code: " << e.code() << endl;
            cout << endl << std::string(90, '-') << "\n\n";
          }

          break;
        }

        case 69: {

          //regex.replace

          std::string str = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];
          std::string regstr = parser(v.Args[1], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];
          std::string replace_with = parser(v.Args[2], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          try {
            std::regex reg(regstr);

            std::string result = std::regex_replace(str, reg, replace_with);

            Action resultJ = strPlaceholder;

            resultJ.ExpStr[0] = result;

            char* cur = "0";

            for (char i : result) {

              Action indexJ = strPlaceholder;

              indexJ.ExpStr = { std::to_string(i) };

              resultJ.Hash_Values[std::string(cur)] = { indexJ };
              cur = AddC(cur, "1", &cli_params.dump()[0]);
            }

            Action retExp = resultJ;

            if (v.SubCall.size() > 0) {

              Action callProcessParser = v;

              bool isProc = v.SubCall[0].IsProc;

              callProcessParser.Indexes = v.SubCall[0].Indexes;
              callProcessParser.Args = v.SubCall[0].Args;
              callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

              retExp = processParser(resultJ, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

            }

            if (expReturn) {

              Returner ret;
              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = retExp;
              ret.type = "expression";

              return ret;
            }

          } catch (std::regex_error& e) {

            //give information about the warning
            cout << "Warning during interpreting: Invalid Regular Expression: " << regstr << endl;
            cout << "Error description: " << e.what() << endl;
            cout << "Error code: " << e.code() << endl;
            cout << endl << std::string(90, '-') << "\n\n";
          }

          break;
        }
        case 70: {

          //this
          //"this" will be used to access the levels of the hash

          std::string level = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          //if it is negative, return undef
          if ((bool) IsLessC(&level[0], "0")) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = falseyVal;
            ret.type = "expression";

            return ret;
          }

          //convert level std::string to ulonglong
          unsigned long long level_number = std::stoull(level);

          //if the this level is too high, return undef
          if (level_number >= this_vals.size()) {
            Returner ret;

            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = falseyVal;
            ret.type = "expression";

            return ret;
          }

          std::map<std::string, std::vector<Action>> hash_level = this_vals[level_number];

          //force all indexes in the hash to be public
          for (std::map<std::string, std::vector<Action>>::iterator it = hash_level.begin(); it != hash_level.end(); ++it)
            hash_level[it->first][0].Access = "public";

          Returner ret;

          std::vector<std::string> retNo;

          Action hashPlaceholder = hashVal;
          hashPlaceholder.Hash_Values = hash_level;

          for (auto it = v.Indexes.begin(); it != v.Indexes.end(); ++it) {
            std::string index = parser(*it, cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];
            hashPlaceholder = hashPlaceholder.Hash_Values[index][0];
          }

          Action retExp = hashPlaceholder;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            retExp = processParser(hashPlaceholder, callProcessParser, cli_params, &vars, this_vals, isProc, dir).exp;

          }

          if (expReturn) {

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = retExp;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 77: {

          //exec

          std::string cmd = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0]; //get the command
          std::string put_stdin = parser(v.Args[1], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0]; //get the stdin

          char* cmdC = &cmd[0];
          char* stdinC = &put_stdin[0];

          //get stdout from the cmd
          char* get_stdout = ExecCmd(cmdC, stdinC, &dir[0]);

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            Action stdout_str = strPlaceholder;

            stdout_str.ExpStr = { std::string(get_stdout) };

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = stdout_str;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 15: {

          //read

          std::string output = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp.ExpStr[0];

          cout << output;

          std::string val;
          cin >> val;

          Action omm_str = ommtypes::to_string(val);

          if (expReturn) {

            Returner ret;
            std::vector<std::string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = omm_str;
            ret.type = "expression";

            return ret;
          }

          break;
        }
        case 19: {

          //typeof

          Returner parsed = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir);

          Action exp = parsed.exp;
          Action stringval = strPlaceholder;

          stringval.ExpStr = { exp.Type };

          std::vector<std::string> noRet;

          if (expReturn) return Returner{ noRet, vars, stringval, "expression" };
          break;
        }
        case 26: {

          //ascii

          Action parsed = parser(v.Args[0], cli_params, vars, false, true, this_vals, dir).exp;

          std::vector<std::string> returnNone;

          if (parsed.Type != "string" && expReturn) return Returner{ returnNone, vars, falseyVal, "expression" };
          else {
            std::string val = parsed.ExpStr[0];
            int first = (int) val[0];

            if (expReturn) {

              Action ascVal = val1;

              ascVal.ExpStr[0] = std::to_string(first);

              return Returner{returnNone, vars, ascVal, "expression"};
            }
          }
        }
        break;
        //////////////////////////

        //assignment operators
        case 4343: {

          //++

          std::string name = v.Name;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              std::vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

              Action val = add(_val, val1, cli_params, this_vals, dir);

              nVar = Variable{
                vars[name].type,
                name,
                { val }
              };
            } else nVar = {
                "local",
                name,
                { val1 }
              };
          }

          vars[name] = nVar;

          if (expReturn) {
            Returner ret;

            std::vector<std::string> retNo;

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

            std::string name = v.Name;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = subtract(_val, val1, cli_params, this_vals, dir);

                nVar = Variable{
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = add(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = subtract(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = multiply(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = divide(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = exponentiate(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

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

            std::string name = v.Name;

            std::vector<Action> __inc = v.ExpAct;
            Action _inc = parser(__inc, cli_params, vars, false, true, this_vals, dir).exp;

            Variable nVar;

            if (vars.find(name) != vars.end()) {

              if (vars[name].type != "dynamic" && vars[name].type != "process") {

                std::vector<Action> __val = vars[name].value;

                Action _val = parser(__val, cli_params, vars, false, true, this_vals, dir).exp;

                Action val = modulo(_val, _inc, cli_params, this_vals, dir);

                nVar = {
                  vars[name].type,
                  name,
                  { val }
                };
              } else nVar = {
                  "local",
                  name,
                  { val1 }
                };
            }

            vars[name] = nVar;

            if (expReturn) {
              Returner ret;

              std::vector<std::string> retNo;

              ret.value = retNo;
              ret.variables = vars;
              ret.exp = _inc;
              ret.type = "expression";

              return ret;
            }
          }
          break;
      }
    }

    std::vector<std::string> returnNone;

    Returner ret;

    ret.value = returnNone;
    ret.variables = vars;
    ret.exp = falseyVal;
    ret.type = "none";

    return ret;
  }
}

#endif
