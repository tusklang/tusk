#include <deque>
#include <vector>
#include "json.hpp"
#include "bind.h"
#include "structs.h"
using namespace std;
using json = nlohmann::json;

Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line);

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

      exp[gen].insert(exp[gen].begin() + spec, evaled.begin(), evaled.end());
    }

    for (int i = 0; i < exp.size(); i++) {
      for (int o = 0; o < exp[i].size(); o++) {

        if (exp[i][o].dump().substr(1, exp[i][o].dump().length() - 2).rfind("$", 0) == 0) {

          json var = vars[exp[i][o].dump().substr(1, exp[i][o].dump().length() - 2)];

          if (var["type"].dump() == "\"process\"") {
            cout << "There Was An Error On Line " << line << ": You cannot have a process in an expression without the '#' keyword"
            << "\n\n" << ((string) var["name"]).substr(1) << endl << "^ <-- Expected '#' here" << endl;

            Kill();
          }

          if (var["value"][0][0].dump() != "null") {
            exp[i].erase(exp[i].begin() + o, exp[i].begin() + o + 1);

            for (int j = 0; j < var["value"][0].size(); j++)
              exp[i].insert(exp[i].begin() + j + o, var["value"][0][j]);
          } else {
            exp[i].erase(exp[i].begin() + o, exp[i].begin() + o + 1);

            for (int j = 0; j < parser(var["valueActs"], calc_params, vars, dir, false, line).exp[0].size(); j++)
              exp[i].insert(exp[i].begin() + j + o, parser(var["valueActs"], calc_params, vars, dir, false, line).exp[0][j]);
          }
        }
      }
    }
    
    //TODO: for each operation, maybe re-program into c++ or even better, fortran

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
