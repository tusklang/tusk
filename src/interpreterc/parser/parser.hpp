#include <iostream>
#include <vector>
#include <deque>
#include "json.hpp"
#include "bind.h"
#include "indexes.hpp"
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

json math(json exp, const json calc_params, json vars, const string dir, int line);

Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line) {

  json expStr = "[]"_json;

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

            //dynamic

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
      }
    } catch (int e) {
      cout << "There Was An Unidentified Error On Line " << line << endl;
      Kill();
    }
  }

  vector<string> returnNone;

  return Returner{ returnNone, vars, math(expStr, calc_params, vars, dir, line), "none" };
}

bool expContain(json exp, string check) {

  for (int i = 0; i < exp.size(); i++) for (int o = 0; o < exp[i].size(); o++) if (exp[i][o] == check) return true;

  return false;
}

tuple<int, int> expIndex(json exp, string check) {
  for (int i = 0; i < exp.size(); i++) for (int o = 0; o < exp[i].size(); o++) if (exp[i][o] == check) return { i, o };

  return { -1, -1 };
}

json math(json exp, const json calc_params, json vars, const string dir, int line) {

  if (exp[0][0] == "true" || exp[0][0] == "false") return exp;
  else {

    while (expContain(exp, "(") && expContain(exp, ")")) {

      int gen, spec;

      tie(gen, spec) = expIndex(exp, "(");

      deque<string> parenExp;

      json part = exp[gen];

      int pCnt = 0;

      for (int i = spec; i < part.size(); i++) {

        if (part[i] == "(") pCnt++;
        if (part[i] == ")") pCnt--;

        parenExp.push_back(part[i]);

        if (pCnt == 0) break;
      }

      parenExp.pop_front();
      parenExp.pop_back();

      vector<string> parenExpVect;

      while (!parenExp.empty()) {

        parenExpVect.push_back(parenExp.front());
        parenExp.pop_front();
      }

      json parenExpJSON_ = parenExpVect;
      string parenExpStr_ = parenExpJSON_.dump();
      char* parenExpStr = &parenExpStr_[0];

      string actions_ = Cactions(parenExpStr);

      json actions = json::parse(actions_);

      Returner parsed = parser(actions, calc_params, vars, dir, false, line);

      json evaled = parsed.exp[0];

      exp[gen].erase(exp[gen].begin() + spec, exp[gen].begin() + parenExpJSON_.size() + 2);
      exp[gen].insert(exp[gen].begin() + spec, evaled[0].begin(), evaled[0].end());
    }

    //for each operation, maybe re-program into c++

    while (expContain(exp, "^")) {

      int gen, spec;

      tie(gen, spec) = expIndex(exp, "^");

      string num1 = exp[gen][spec - 1]
      , num2 = exp[gen][spec + 1];

      string cp = calc_params.dump();

      char* val = Exponentiate(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

      exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

      exp[gen].insert(exp[gen].begin() + spec - 1, val);
    }

    while (expContain(exp, "*") || expContain(exp, "/")) {

      int multg, mults, divg, divs;

      tie(multg, mults) = expIndex(exp, "*");
      tie(divg, divs) = expIndex(exp, "-");

      if (multg > divg || mults > divs || divs == -1) {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "*");

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Multiply(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      } else {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "/");

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Division(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      }
    }

    while (expContain(exp, "%")) {

      int gen, spec;

      tie(gen, spec) = expIndex(exp, "%");

      string num1 = exp[gen][spec - 1]
      , num2 = exp[gen][spec + 1];

      string cp = calc_params.dump();

      char* val = Modulo(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

      exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

      exp[gen].insert(exp[gen].begin() + spec - 1, val);
    }

    while (expContain(exp, "+") || expContain(exp, "-")) {

      int addg, adds, subg, subs;

      tie(addg, adds) = expIndex(exp, "+");
      tie(subg, subs) = expIndex(exp, "-");

      if (addg > subg || adds > subs || subs == -1) {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "+");

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Add(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      } else {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "-");

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Subtract(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      }
    }

    return exp;
  }
}
