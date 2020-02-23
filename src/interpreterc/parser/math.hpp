#include <tuple>
#include <iostream>
#include <deque>
#include <vector>
#include "json.hpp"
#include "operations/add.hpp"
#include "operations/subtract.hpp"
#include "operations/multiply.hpp"
#include "operations/divide.hpp"
#include "operations/exponent.hpp"
#include "operations/modulo.hpp"
using namespace std;
using json = nlohmann::json;

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

      //maybe switch to a deque
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

      json parenExpJSON_ = parenExpVect
      , parenExpJSON = json::parse("[" + parenExpJSON_.dump() + "]");

      json evaled = math(parenExpJSON, calc_params, vars, dir, line);

      exp[gen] = exp[gen].erase(evaled.begin() + spec, evaled.begin() + parenExp.size());

      cout << exp << endl;
      //exp[gen] = exp[gen].insert(spec, evaled[0][0])
    }

    return exp;
  }
}
