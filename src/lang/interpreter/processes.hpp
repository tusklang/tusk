#ifndef PROCESSES_HPP_
#define PROCESSES_HPP_

#include <vector>
#include <deque>
#include <map>

#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
using namespace std;

Returner parser(const vector<Action> actions, const json cli_params, map<string, Variable> vars, const bool groupReturn, const bool expReturn, deque<map<string, vector<Action>>> this_vals, string dir);

Returner processParser(Action var, const Action v, const json cli_params, map<string, Variable>* vars, deque<map<string, vector<Action>>> this_vals, bool isProc, string dir) {

  vector<string> noRet;

  if (v.Indexes.size() == 0 && v.Args.size() == 0) return Returner{ noRet, *vars, var, "expression" };

  Returner fparsed = Returner{ noRet, *vars, falseyVal, "expression" };

  Returner parsed = fparsed;

  deque<map<string, vector<Action>>> send_this = this_vals;

  for (auto it = v.Indexes.begin(); it != v.Indexes.end(); ++it) {

    string index = parser(*it, cli_params, *vars, false, true, this_vals, dir).exp.ExpStr[0];

    if (var.Hash_Values.find(index) == var.Hash_Values.end() || (islower(index[0]) && var.Hash_Values[index][0].Access != "public")) {
      parsed = fparsed;
      return Returner{ noRet, *vars, falseyVal, "expression" };
    }

    send_this.push_front(var.Hash_Values);
    var = parser(var.Hash_Values[index], cli_params, *vars, false, true, this_vals, dir).exp;
  }

  if (!isProc) return Returner{ noRet, *vars, var, "expression" };

  if (var.Type != "process") {
    parsed = fparsed;
    return Returner{ noRet, *vars, falseyVal, "expression" };
  }

  vector<string> params = var.Params;
  vector<vector<Action>> args = v.Args;

  if (params.size() != args.size()) {
    parsed = fparsed;
    return Returner{ noRet, *vars, falseyVal, "expression" };
  }

  map<string, Variable> sendVars = *vars;

  for (int o = 0; o < params.size() || o < args.size(); o++) {

    Variable cur = Variable{
      "local",
      params[o],
      { parser(args[o], cli_params, *vars, false, true, this_vals, dir).exp }
    };

    sendVars[params[o]] = cur;
  }

  parsed = parser(var.ExpAct, cli_params, sendVars, true, false, send_this, dir);

  map<string, Variable> pVars = parsed.variables;

  //filter the variables that are not global
  for (pair<string, Variable> o : pVars)
    if (o.second.type == "global" || o.second.type == "process" || (*vars).find(o.second.name) != (*vars).end())
      (*vars)[o.first] = o.second;

  for (SubCaller it : v.SubCall) {

    Action curVar = parsed.exp;

    curVar.Indexes = it.Indexes;
    curVar.Args = it.Args;

    parsed = processParser(curVar, curVar, cli_params, vars, this_vals, it.IsProc, dir);
  }

  Action val = parsed.exp;

  return Returner{ noRet, *vars, val, "expression" };
}

#endif
