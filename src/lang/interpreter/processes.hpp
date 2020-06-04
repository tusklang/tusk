#ifndef OMM_PROCESSES_HPP_
#define OMM_PROCESSES_HPP_

#include <vector>
#include <deque>
#include <map>
#include <algorithm>

#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"

namespace omm {

  Returner parser(const std::vector<Action> actions, const json cli_params, std::map<std::string, Variable> vars, const bool groupReturn, const bool expReturn, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir);

  bool vector_indexes_inc(std::vector<std::string> vec, std::string str) {

    for (std::string v : vec)

      //see if v includes str
      if (v.find(str) != std::string::npos) return true;

    return false;
  }

  Returner processParser(Action var, const Action v, const json cli_params, std::map<std::string, Variable>* vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, bool isProc, std::string dir) {

    std::vector<std::string> noRet;

    Returner fparsed = Returner{ noRet, *vars, falseyVal, "expression" };

    Returner parsed = fparsed;

    std::deque<std::map<std::string, std::vector<Action>>> send_this = this_vals;

    for (auto it = v.Indexes.begin(); it != v.Indexes.end(); ++it) {

      std::string index = parser(*it, cli_params, *vars, false, true, this_vals, dir).exp.ExpStr[0];

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
      return Returner{ noRet, *vars, var, "expression" };
    }

    std::vector<std::string> params = var.Params;
    std::vector<std::vector<Action>> args = v.Args;

    if (params.size() != args.size() && !vector_indexes_inc(params, "pargv")) {
      parsed = fparsed;
      return Returner{ noRet, *vars, falseyVal, "expression" };
    }

    std::map<std::string, Variable> sendVars = *vars;

    for (int o = 0; o < params.size() || o < args.size(); o++) {

      //if it starts with pargv
      if (params[o].rfind("$pargv.", 0) == 0) {

        std::string varname = "$" + params[o].substr(std::string("$pargv.").length());

        //convert the rest of the args into an array and store it in the pargv variable
        std::map<std::string, std::vector<Action>> pargv;

        for (unsigned long long cur = 0; o < args.size(); ++o, ++cur)
          pargv[to_string(cur)] = { omm::parser(args[o], cli_params, *vars, false, true, this_vals, dir).exp };

        Action arg = arrayVal;
        arg.Hash_Values = pargv;

        sendVars[varname] = Variable{
          "pargv",
          varname,
          { arg },
          [](Action v, json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) -> Returner { return Returner{}; }
        };

        break;
      }

      Variable cur = Variable{
        "argument",
        params[o],
        { omm::parser(args[o], cli_params, *vars, false, true, this_vals, dir).exp },
        [](Action v, json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) -> Returner { return Returner{}; }
      };

      sendVars[params[o]] = cur;
    }

    parsed = parser(var.ExpAct, cli_params, sendVars, true, true, send_this, dir);

    std::map<std::string, Variable> pVars = parsed.variables;

    //filter the variables that are not global
    for (std::pair<std::string, Variable> o : pVars)
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

}

#endif
