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

//file i/o
#include "../files/readfile.hpp"
#include "../files/writefile.h"
#include "../files/delete.h"
#include "../files/isDir.hpp"
#include "../files/isFile.hpp"

//operations
#include "operations/add/add.hpp"
#include "operations/divide/divide.hpp"
#include "operations/exponentiate/exponentiate.hpp"
#include "operations/modulo/modulo.hpp"
#include "operations/multiply/multiply.hpp"
#include "operations/subtract/subtract.hpp"

using namespace std;
using json = nlohmann::json;
using ulong = unsigned long;

Returner parser(const vector<Action> actions, const json cli_params, map<string, Variable> vars, const bool groupReturn, const bool expReturn, deque<map<string, vector<Action>>> this_vals) {

  //loop through every action
  for (Action v : actions) {

    //get current action id
    int cur = v.ID;

    switch (cur) {
      case 1: {

          //local

          string name = v.Name;

          vector<Action> acts = v.ExpAct
          , parsed = { parser(acts, cli_params, vars, false, true, this_vals).exp };

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

          string name = v.Name;

          vector<Action> acts = v.ExpAct;

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

          Returner cond = parser(v.Condition[0].Condition, cli_params, vars, false, true, this_vals);

          //while the alt statement should continue
          while (isTruthy(cond.exp)) {

            //going back to the first block when it reached the last block
            if (o >= v.Condition.size()) o = 0;

            Returner parsed = parser(v.Condition[o].Actions, cli_params, vars, true, false, this_vals);

            map<string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (pair<string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                vars[o.first] = o.second;

            if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
            if (parsed.type == "skip") continue;
            if (parsed.type == "break") break;

            cond = parser(v.Condition[o].Condition, cli_params, vars, false, true, this_vals);
            o++;
          }
        }
        break;
      case 4: {

          //global

          string name = v.Name;

          vector<Action> acts = v.ExpAct
          , parsed = { parser(acts, cli_params, vars, false, true, this_vals).exp };

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

          Action _val = parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp;

          log_format(_val, cli_params, vars, 2, "log");
        }
        break;
      case 6: {

          //print

          Action _val = parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp;

          log_format(_val, cli_params, vars, 2, "print");
        }
        break;
      case 8: {

          //expressionIndex

          Action val = parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp;

          Action index = indexesCalc(val.Hash_Values, v.Indexes, cli_params, vars, this_vals);

          if (expReturn) {
            vector<string> returnNone;

            return Returner{ returnNone, vars, index, "expression" };
          }
        }
        break;
      case 9: {

          //group

          vector<Action> acts = v.ExpAct;

          Returner parsed = parser(acts, cli_params, vars, false, false, this_vals);

          map<string, Variable> pVars = parsed.variables;

          //filter the variables that are not global
          for (pair<string, Variable> o : pVars)
            if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
              vars[o.first] = o.second;

          if (groupReturn || expReturn) return Returner{ parsed.value, vars, parsed.exp, parsed.type };
        }
        break;
      case 10: {

          //process

          string name = v.Name;

          if (name != "") {
            Variable nVar = Variable{
              "process",
              name,
              { v }
            };

            vars[name] = nVar;
          }

          if (expReturn) {
            vector<string> noRet;

            return Returner{ noRet, vars, v, "expression" };
          }
        }
        break;
      case 11: {

          //# (call process)

          string name = v.Name;

          Returner parsed;

          vector<string> noRet;

          Returner fparsed = Returner{ noRet, vars, falseyVal, "expression" };

          parsed = fparsed;

          if (vars.find(name) == vars.end()) goto stopIndexing_processes;
          else {

            Action var = vars[name].value[0];

            parsed = processParser(var, v, cli_params, &vars, this_vals, true);
          }

          stopIndexing_processes:
          if (expReturn) {

            Action val = parsed.exp;

            vector<string> noRet;

            return Returner{ noRet, vars, val, "expression" };
          }
        }
        break;
      case 12: {

          //return

          vector<string> noRet;

          return Returner{ noRet, vars, parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp, "return" };
        }
        break;
      case 13: {

          //conditional

          for (int o = 0; o < v.Condition.size(); o++) {

            Action val = parser(v.Condition[o].Condition, cli_params, vars, false, true, this_vals).exp;

            if (isTruthy(val)) {

              Returner parsed = parser(v.Condition[o].Actions, cli_params, vars, true, false, this_vals);

              map<string, Variable> pVars = parsed.variables;

              //filter the variables that are not global
              for (pair<string, Variable> o : pVars)
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

          vector<vector<Action>> files = v.Value; //get all actionized files imported

          //loop through actionized files
          for (vector<Action> it : files) {

            Returner parsed = parser(it, cli_params, vars, true, false, this_vals);

            map<string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (pair<string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process")
                vars[v.Name + "." + o.first.substr(1)] = o.second;
          }
        }
        break;
      case 15: {

          //read

          string in;

          cout << ((string) parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp.ExpStr[0]) << " ";

          cin >> in;

          if (expReturn) {
            vector<string> retNo;

            Action expRet = strPlaceholder;

            expRet.ExpStr[0] = in;

            for (char c : in) {

              Action cPlaceholder = strPlaceholder;

              cPlaceholder.ExpStr[0] = to_string(c);
            }

            return Returner{ retNo, vars, expRet, "expression" };
          }
        }
        break;
      case 16: {

          //break

          Returner ret;

          vector<string> returnNone;
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

          vector<string> returnNone;
          Action expNone;

          ret.value = returnNone;
          ret.variables = vars;
          ret.exp = expNone;
          ret.type = "skip";

          return ret;
        }
        break;
      case 19: {

          //typeof

          Returner parsed = parser(v.ExpAct, cli_params, vars, false, true, this_vals);

          Action exp = parsed.exp;
          Action stringval = strPlaceholder;

          stringval.ExpStr = { exp.Type };

          vector<string> noRet;

          if (expReturn) return Returner{ noRet, vars, stringval, "expression" };
        }
        break;
      case 21: {

          //loop

          vector<Action> cond = v.Condition[0].Condition
          , acts = v.Condition[0].Actions;

          Returner parsed;

          Action condP = parser(cond, cli_params, vars, false, true, this_vals).exp;

          while (isTruthy(condP)) {

            parsed = parser(acts, cli_params, vars, true, false, this_vals);

            map<string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (pair<string, Variable> o : pVars)
              if (o.second.type == "global" || o.second.type == "process" || vars.find(o.second.name) != vars.end())
                vars[o.first] = o.second;

            if (parsed.type == "return") return Returner{ parsed.value, vars, parsed.exp, "return" };
            if (parsed.type == "skip") continue;
            if (parsed.type == "break") break;

            condP = parser(cond, cli_params, vars, false, true, this_vals).exp;
          }

        }
        break;
      case 22: {

          //hash

          if (expReturn) {

            vector<string> returnNone;

            bool isMutable = v.IsMutable;

            Action val = v;

            if (!isMutable) {

              for (pair<string, vector<Action>> it : v.Hash_Values) {

                string paramCount = "";
                Action exp = parser(it.second, cli_params, vars, false, true, this_vals).exp;

                val.Hash_Values[it.first] = { exp };
              }
            }

            return Returner{ returnNone, vars, val, "expression" };
          }
        }
        break;
      case 23: {

          //hashIndex

          map<string, vector<Action>> val = v.Hash_Values;

          Action index = indexesCalc(val, v.Indexes, cli_params, vars, this_vals);

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

            bool isMutable = v.IsMutable;

            Action val = v;

            if (!isMutable) {

              char* index = "0";

              for (pair<string, vector<Action>> o : v.Hash_Values) {

                if (val.Hash_Values.find(index) == val.Hash_Values.end()) {
                  index = AddC(index, "1", &cli_params.dump()[0]);
                  continue;
                }

                string paramCount = "";
                Action exp = parser(o.second, cli_params, vars, false, true, this_vals).exp;

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

          map<string, vector<Action>> val = v.Hash_Values;
          Action index = indexesCalc(val, v.Indexes, cli_params, vars, this_vals);

          if (expReturn) {
            vector<string> returnNone;

            return Returner{ returnNone, vars, index, "expression" };
          }
        }
        break;
      case 26: {

          //ascii

          Action parsed = parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp;

          vector<string> returnNone;

          if (parsed.Type != "string" && expReturn) return Returner{ returnNone, vars, falseyVal, "expression" };
          else {
            string val = parsed.ExpStr[0];
            int first = (int) val[0];

            if (expReturn) {

              Action ascVal = val1;

              ascVal.ExpStr[0] = to_string(first);

              return Returner{returnNone, vars, ascVal, "expression"};
            }
          }
        }
        break;
      case 28: {

          //let

          string name = v.Name;

          vector<Action> acts = v.ExpAct;

          vector<Action> parsed = { parser(acts, cli_params, vars, false, true, this_vals).exp };

          Variable nVar;

          Variable var = vars[name];
          vector<string> indexes;

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

            for (vector<Action> it : v.Indexes) {

              string varP = parser(it, cli_params, vars, false, true, this_vals).exp.ExpStr[0];

              map = &(map->Hash_Values[varP][0]);
            }

            *map = parsed[0];

          }
        }
        break;
      case 32: {

          //add

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = add(first, second, cli_params, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = subtract(first, second, cli_params, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = multiply(first, second, cli_params, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = divide(first, second, cli_params, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = exponentiate(first, second, cli_params, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = modulo(first, second, cli_params, this_vals);

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

          if (expReturn) return Returner{ noRet, vars, v, "expression" };
        }
        break;
      case 39: {

          //number

          vector<string> noRet;

          if (expReturn) return Returner{ noRet, vars, v, "expression" };
        }
        break;
      case 40: {

          //boolean

          vector<string> noRet;

          if (expReturn) return Returner{ noRet, vars, v, "expression" };
        }
        break;
      case 41: {

          //falsey

          vector<string> noRet;

          if (expReturn) return Returner{ noRet, vars, v, "expression" };
        }
        break;
      case 42: {

          //none

          vector<string> noRet;

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

            val = parser({ var }, cli_params, vars, false, true, this_vals).exp;
          }

          vector<string> noRet;

          if (expReturn) return Returner{ noRet, vars, val, "expression" };
        }
        break;
      case 46: {

          //variableIndex

          Returner parsedVal = parser(v.ExpAct, cli_params, vars, false, true, this_vals);

          Action index = indexesCalc(parsedVal.exp.Hash_Values, v.Indexes, cli_params, vars, this_vals);

          if (expReturn) {
            vector<string> returnNone;

            return Returner{ returnNone, vars, index, "expression" };
          }
        }
        break;
      case 47: {

          //equals

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action val = equals(
            first,
            second,
            cli_params,
            vars,
            this_vals
          );

          if (first.Type != second.Type) val = falseRet;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        Action val = equals(
          first,
          second,
          cli_params,
          vars,
          this_vals
        );

        val = val.ExpStr[0] == "true" ? falseRet : trueRet;
        if (first.Type != second.Type) val = trueRet;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        Action val = isGreater(
          first,
          second,
          cli_params
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
        break;
      }
      case 50: {

        //less

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        Action val = isLess(
          first,
          second,
          cli_params
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
        break;
      }
      case 51: {

        //greaterOrEqual

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        Action val = isLess(
          first,
          second,
          cli_params
        );

        val = val.ExpStr[0] == "true" ? falseRet : trueRet;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        Action val = isGreater(
          first,
          second,
          cli_params
        );

        val = val.ExpStr[0] == "true" ? falseRet : trueRet;
        if (first.Type != second.Type) val = trueRet;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

          Action val = parser(v.Second, cli_params, vars, false, true, this_vals).exp
          , retval;

          string expstr = val.ExpStr[0];

          if (expstr == "false" || val.Type == "falsey") retval = trueRet;
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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action retval;

          if ((v.Degree).size() == 0) retval = similarity(first, second, zero, cli_params, vars, this_vals);
          else retval = similarity(first, second, parser(v.Degree, cli_params, vars, false, true, this_vals).exp, cli_params, vars, this_vals);

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

          Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
          , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

          Action retval;

          if ((v.Degree).size() == 0) retval = strictSimilarity(first, second, zero, cli_params, vars, this_vals);
          else retval = strictSimilarity(first, second, parser(v.Degree, cli_params, vars, false, true, this_vals).exp, cli_params, vars, this_vals);

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
      case 71: {

        //or

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

        Action first = parser(v.First, cli_params, vars, false, true, this_vals).exp
        , second = parser(v.Second, cli_params, vars, false, true, this_vals).exp;

        if (expReturn) {
          Returner ret;

          vector<string> retNo;

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

          string name = v.Name;

          Returner parsed;

          vector<string> noRet;

          Returner fparsed = Returner{ noRet, vars, falseyVal, "none" };

          parsed = fparsed;

          if (vars.find(name) == vars.end()) goto stopIndexing_threads;
          else {

            Action var = vars[name].value[0];

            for (vector<Action> it : v.Indexes) {

              string index = parser(it, cli_params, vars, false, true, this_vals).exp.ExpStr[0];

              if (var.Hash_Values.find(index) == var.Hash_Values.end()) {
                parsed = fparsed;
                goto stopIndexing_threads;
              }

              var = parser(var.Hash_Values[index], cli_params, vars, false, true, this_vals).exp;
            }

            if (var.Type != "process") {
              parsed = fparsed;
              goto stopIndexing_threads;
            }

            vector<string> params = var.Params;
            vector<vector<Action>> args = v.Args;

            if (params.size() != args.size()) {
              parsed = fparsed;
              goto stopIndexing_threads;
            }

            map<string, Variable> sendVars = vars;

            for (int o = 0; o < params.size() || o < args.size(); o++) {

              Variable cur = Variable{
                "local",
                params[o],
                { parser(args[o], cli_params, vars, false, true, this_vals).exp }
              };

              sendVars[params[o]] = cur;
            }

            thread _(parser, var.ExpAct, cli_params, sendVars, true, false, this_vals);

            _.detach();
          }

          stopIndexing_threads:
          if (expReturn) {

            Action val = parsed.exp;

            vector<string> noRet;

            return Returner{ noRet, vars, val, "expression" };
          }
        }
        break;
      case 57: {

          //wait

          Action amt = parser(v.ExpAct, &cli_params.dump()[0], vars, false, true, this_vals).exp;

          if (IsLessC(&(amt.ExpStr[0])[0], "4294967296")) Sleep((ulong) atoi(&(amt.ExpStr[0])[0]));
          else {
            for (char* o = "0"; (bool) IsLessC(o, &(amt.ExpStr[0])[0]); o = AddC(o, "4294967296", &cli_params.dump()[0])) {

              char* subtracted = SubtractC(&(amt.ExpStr[0])[0], o, &cli_params.dump()[0]);

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

            Action cur = parser(v.ExpAct, cli_params, vars, false, true, this_vals).exp;
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

          vector<string> putterVars = v.ExpStr;
          string var1 = putterVars[0]
          , var2 = putterVars[1];

          //parse the iterator value
          map<string, vector<Action>> iterator = parser(v.First /* v.First is where the iterator is stored */, cli_params, vars, false, true, this_vals).exp.Hash_Values;

          iterator.erase("falsey");

          for (pair<string, vector<Action>> it : iterator) {
            map<string, Variable> sendVars = vars;

            Action key = strPlaceholder;

            sendVars[var1] = Variable{
              "local",
              var1,
              { key }
            };
            sendVars[var2] = Variable{
              "local",
              var2,
              { parser(it.second, cli_params, vars, false, true, this_vals).exp }
            };

            Returner parsed = parser(v.ExpAct, cli_params, sendVars, true, false, this_vals);

            map<string, Variable> pVars = parsed.variables;

            //filter the variables that are not global
            for (pair<string, Variable> o : pVars)
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

        string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        smatch match;

        //see if the filename is absolute
        regex pat("^[a-zA-Z]:");
        bool isOnDrive = regex_search(filename, match, pat);

        string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

        if (!isFile(nDir + filename) && expReturn) {

          Action returner = falseyVal;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

          }

          if (expReturn) {

            Returner ret;
            vector<string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returner;
            ret.type = "expression";

            return ret;
          }

        } else {
          string content = readfile(&(nDir + filename)[0]);

          if (expReturn) {
            Returner ret;

            vector<string> retNo;

            Action contentJ = strPlaceholder;

            contentJ.ExpStr = {content};

            //make the hash values of the string
            for (ulong o = 0; o < content.length(); o++) {
              Action curChar = strPlaceholder;

              curChar.ExpStr = {
                string(1, content[o])
              };

              contentJ.Hash_Values[to_string(o)] = { curChar };
            }

            Action returner = contentJ;

            if (v.SubCall.size() > 0) {

              Action callProcessParser = v;

              bool isProc = v.SubCall[0].IsProc;

              callProcessParser.Indexes = v.SubCall[0].Indexes;
              callProcessParser.Args = v.SubCall[0].Args;
              callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

              returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

            }

            if (expReturn) {

              Returner ret;
              vector<string> retNo;

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
        string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];
        Action content = parser(v.Args[1], cli_params, vars, false, true, this_vals).exp;

        smatch match;

        //see if the filename is absolute
        regex pat("^[a-zA-Z]:");
        bool isOnDrive = regex_search(filename, match, pat);

        string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

        if (content.Type == "falsey") {

          deletefile(&(nDir + filename)[0]);

        } else {

          string contentstr = content.ExpStr[0];
          writefile(&(nDir + filename)[0], &contentstr[0]);
        }

        Action returner = content;

        if (v.SubCall.size() > 0) {

          Action callProcessParser = v;

          bool isProc = v.SubCall[0].IsProc;

          callProcessParser.Indexes = v.SubCall[0].Indexes;
          callProcessParser.Args = v.SubCall[0].Args;
          callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

          returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

        }

        if (expReturn) {

          Returner ret;
          vector<string> retNo;

          ret.value = retNo;
          ret.variables = vars;
          ret.exp = returner;
          ret.type = "expression";

          return ret;
        }
        break;
      }
      case 62: {

        //files.exists

        //written as file.exists(dir)

        string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        smatch match;

        //see if the filename is absolute
        regex pat("^[a-zA-Z]:");
        bool isOnDrive = regex_search(filename, match, pat);

        string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

        //if it is not a directory and not a file, it does not exist
        bool exists = !(!isDir(nDir + filename) && !isFile(nDir + filename));

        Action returner = exists ? trueRet : falseRet;

        if (v.SubCall.size() > 0) {

          Action callProcessParser = v;

          bool isProc = v.SubCall[0].IsProc;

          callProcessParser.Indexes = v.SubCall[0].Indexes;
          callProcessParser.Args = v.SubCall[0].Args;
          callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

          returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

        }

        if (expReturn) {

          Returner ret;
          vector<string> retNo;

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

        //written as file.isFile(dir)

        string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        smatch match;

        //see if the filename is absolute
        regex pat("^[a-zA-Z]:");
        bool isOnDrive = regex_search(filename, match, pat);

        string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

        bool isFileVal = isFile(nDir + filename);

        Action returner = isFileVal ? trueRet : falseRet;

        if (v.SubCall.size() > 0) {

          Action callProcessParser = v;

          bool isProc = v.SubCall[0].IsProc;

          callProcessParser.Indexes = v.SubCall[0].Indexes;
          callProcessParser.Args = v.SubCall[0].Args;
          callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

          returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

        }

        if (expReturn) {

          Returner ret;
          vector<string> retNo;

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

        //written as file.isDir(dir)

        string filename = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        smatch match;

        //see if the filename is absolute
        regex pat("^[a-zA-Z]:");
        bool isOnDrive = regex_search(filename, match, pat);

        string nDir = isOnDrive ? "" : cli_params["Files"]["DIR"];

        bool isDirVal = isDir(nDir + filename);

        Action returner = isDirVal ? trueRet : falseRet;

        if (v.SubCall.size() > 0) {

          Action callProcessParser = v;

          bool isProc = v.SubCall[0].IsProc;

          callProcessParser.Indexes = v.SubCall[0].Indexes;
          callProcessParser.Args = v.SubCall[0].Args;
          callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

          returner = processParser(returner, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

        }

        if (expReturn) {

          Returner ret;
          vector<string> retNo;

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

        string str = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];
        string regstr = parser(v.Args[1], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        try {
          regex reg(regstr);

          smatch matcher;

          vector<ulong long> found_indexes;

          //get all matches
          for (auto it = sregex_iterator(str.begin(), str.end(), reg); it != sregex_iterator(); it++) {
            found_indexes.push_back(it->position());
          }

          Action returnerArr = arrayVal;

          char* cur = "0";

          //loop through the indexes found and store them an omm type array
          for (int i : found_indexes) {

            //store the value of the number 1
            Action indexJ = val1;

            indexJ.ExpStr[0] = to_string(i);

            returnerArr.Hash_Values[string(cur)] = { indexJ };
            cur = AddC(cur, "1", &cli_params.dump()[0]);
          }

          Action returnerVal = hashVal;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            returnerVal = processParser(returnerArr, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

          }

          if (expReturn) {

            Returner ret;
            vector<string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = returnerVal;
            ret.type = "expression";

            return ret;
          }

        } catch (regex_error& e) {

          //give information about the warning
          cout << "Warning during interpreting: Invalid Regular Expression: " << regstr << endl;
          cout << "Error description: " << e.what() << endl;
          cout << "Error code: " << e.code() << endl;
          cout << endl << string(90, '-') << "\n\n";
        }

        break;
      }

      case 69: {

        //regex.replace

        string str = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];
        string regstr = parser(v.Args[1], cli_params, vars, false, true, this_vals).exp.ExpStr[0];
        string replace_with = parser(v.Args[2], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        try {
          regex reg(regstr);

          string result = regex_replace(str, reg, replace_with);

          Action resultJ = strPlaceholder;

          resultJ.ExpStr[0] = result;

          char* cur = "0";

          for (char i : result) {

            Action indexJ = strPlaceholder;

            indexJ.ExpStr = { to_string(i) };

            resultJ.Hash_Values[string(cur)] = { indexJ };
            cur = AddC(cur, "1", &cli_params.dump()[0]);
          }

          Action retExp = resultJ;

          if (v.SubCall.size() > 0) {

            Action callProcessParser = v;

            bool isProc = v.SubCall[0].IsProc;

            callProcessParser.Indexes = v.SubCall[0].Indexes;
            callProcessParser.Args = v.SubCall[0].Args;
            callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

            retExp = processParser(resultJ, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

          }

          if (expReturn) {

            Returner ret;
            vector<string> retNo;

            ret.value = retNo;
            ret.variables = vars;
            ret.exp = retExp;
            ret.type = "expression";

            return ret;
          }

        } catch (regex_error& e) {

          //give information about the warning
          cout << "Warning during interpreting: Invalid Regular Expression: " << regstr << endl;
          cout << "Error description: " << e.what() << endl;
          cout << "Error code: " << e.code() << endl;
          cout << endl << string(90, '-') << "\n\n";
        }

        break;
      }
      case 70: {

        //this
        //"this" will be used to access the levels of the hash

        string level = parser(v.Args[0], cli_params, vars, false, true, this_vals).exp.ExpStr[0];

        //if it is negative, return undef
        if ((bool) IsLessC(&level[0], "0")) {
          Returner ret;

          vector<string> retNo;

          ret.value = retNo;
          ret.variables = vars;
          ret.exp = falseyVal;
          ret.type = "expression";

          return ret;
        }

        //convert level string to ulonglong
        unsigned long long level_number = stoull(level);

        //if the this level is too high, return undef
        if (level_number >= this_vals.size()) {
          Returner ret;

          vector<string> retNo;

          ret.value = retNo;
          ret.variables = vars;
          ret.exp = falseyVal;
          ret.type = "expression";

          return ret;
        }

        map<string, vector<Action>> hash_level = this_vals[level_number];

        //force all indexes in the hash to be public
        for (map<string, vector<Action>>::iterator it = hash_level.begin(); it != hash_level.end(); ++it)
          hash_level[it->first][0].Access = "public";

        Returner ret;

        vector<string> retNo;

        Action hashPlaceholder = hashVal;
        hashPlaceholder.Hash_Values = hash_level;

        for (auto it = v.Indexes.begin(); it != v.Indexes.end(); ++it) {
          string index = parser(*it, cli_params, vars, false, true, this_vals).exp.ExpStr[0];
          hashPlaceholder = hashPlaceholder.Hash_Values[index][0];
        }

        Action retExp = hashPlaceholder;

        if (v.SubCall.size() > 0) {

          Action callProcessParser = v;

          bool isProc = v.SubCall[0].IsProc;

          callProcessParser.Indexes = v.SubCall[0].Indexes;
          callProcessParser.Args = v.SubCall[0].Args;
          callProcessParser.SubCall.erase(callProcessParser.SubCall.begin());

          retExp = processParser(hashPlaceholder, callProcessParser, cli_params, &vars, this_vals, isProc).exp;

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
      //////////////////////////

      //assignment operators
      case 4343: {

        //++

        string name = v.Name;

        Variable nVar;

        if (vars.find(name) != vars.end()) {

          if (vars[name].type != "dynamic" && vars[name].type != "process") {

            vector<Action> __val = vars[name].value;

            Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

            Action val = add(_val, val1, cli_params, this_vals);

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

          string name = v.Name;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = subtract(_val, val1, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = add(_val, _inc, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = subtract(_val, _inc, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = multiply(_val, _inc, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = divide(_val, _inc, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = exponentiate(_val, _inc, cli_params, this_vals);

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

          string name = v.Name;

          vector<Action> __inc = v.ExpAct;
          Action _inc = parser(__inc, cli_params, vars, false, true, this_vals).exp;

          Variable nVar;

          if (vars.find(name) != vars.end()) {

            if (vars[name].type != "dynamic" && vars[name].type != "process") {

              vector<Action> __val = vars[name].value;

              Action _val = parser(__val, cli_params, vars, false, true, this_vals).exp;

              Action val = modulo(_val, _inc, cli_params, this_vals);

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
