#ifndef MATH_HPP_
#define MATH_HPP_

#include <iostream>
#include <deque>
#include <vector>
#include <string.h>
#include "json.hpp"
#include "bind.h"
#include "structs.h"
using namespace std;
using json = nlohmann::json;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn);

bool expContain(json exp, string check, vector<int> checked_gens) {

  for (int i = 0; i < exp.size(); i++) for (int o = 0; o < exp[i].size(); o++) if (exp[i][o] == check) {
    for (int j : checked_gens)
      if (i == j) goto outer;
    return true;

    outer:;
  }

  return false;
}

bool expContain(json exp, string check) {

  for (int i = 0; i < exp.size(); i++) for (int o = 0; o < exp[i].size(); o++) if (exp[i][o] == check) return true;

  return false;
}

tuple<int, int> expIndex(json exp, string check, vector<int> checked_gens) {

  for (int i = 0; i < exp.size(); i++) for (int o = 0; o < exp[i].size(); o++) if (exp[i][o] == check) {
    for (int j : checked_gens)
      if (i == j) goto outer;
    return { i, o };

    outer:;
  }

  return { -1, -1 };
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

      Returner parsed = parser(actions, calc_params, vars, dir, false, line, true);

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

            for (int j = 0; j < parser(var["valueActs"], calc_params, vars, dir, false, line, true).exp[0].size(); j++)
              exp[i].insert(exp[i].begin() + j + o, parser(var["valueActs"], calc_params, vars, dir, false, line, true).exp[0][j]);
          }
        }
      }
    }

    //TODO: for each operation, maybe re-program into c++ or even better, fortran

    vector<int> checked_gens_exponent;

    while (expContain(exp, "^", checked_gens_exponent)) {

      int gen, spec;

      tie(gen, spec) = expIndex(exp, "^", checked_gens_exponent);

      if (exp[gen].size() < 3) {
        checked_gens_exponent.push_back(gen);
        continue;
      }

      string num1 = exp[gen][spec - 1]
      , num2 = exp[gen][spec + 1];

      string cp = calc_params.dump();

      char* val = Exponentiate(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

      exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

      exp[gen].insert(exp[gen].begin() + spec - 1, val);
    }

    vector<int> checked_gens_multiply, checked_gens_divide;

    while (expContain(exp, "*", checked_gens_multiply) || expContain(exp, "/", checked_gens_divide)) {

      int multg, mults, divg, divs;

      tie(multg, mults) = expIndex(exp, "*");
      tie(divg, divs) = expIndex(exp, "/");

      if (multg > divg || mults > divs || divs == -1) {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "*", checked_gens_multiply);

        if (exp[gen].size() < 3) {
          checked_gens_multiply.push_back(gen);
          continue;
        }

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Multiply(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      } else {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "/", checked_gens_divide);

        if (exp[gen].size() < 3) {
          checked_gens_divide.push_back(gen);
          continue;
        }

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Division(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      }
    }

    vector<int> checked_gens_modulo;

    while (expContain(exp, "%", checked_gens_modulo)) {

      int gen, spec;

      tie(gen, spec) = expIndex(exp, "%", checked_gens_modulo);

      if (exp[gen].size() < 3) {
        checked_gens_modulo.push_back(gen);
        continue;
      }

      string num1 = exp[gen][spec - 1]
      , num2 = exp[gen][spec + 1];

      string cp = calc_params.dump();

      char* val = Modulo(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

      exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

      exp[gen].insert(exp[gen].begin() + spec - 1, val);
    }

    vector<int> checked_gens_add, checked_gens_subtract;

    while (expContain(exp, "+", checked_gens_add) || expContain(exp, "-", checked_gens_subtract)) {

      int addg, adds, subg, subs;

      tie(addg, adds) = expIndex(exp, "+", checked_gens_add);
      tie(subg, subs) = expIndex(exp, "-", checked_gens_subtract);

      if (addg > subg || adds > subs || subs == -1) {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "+");

        if (exp[gen].size() < 3) {
          checked_gens_add.push_back(gen);
          continue;
        }

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Add(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      } else {
        int gen, spec;

        tie(gen, spec) = expIndex(exp, "-");

        if (exp[gen].size() < 3) {
          checked_gens_subtract.push_back(gen);
          continue;
        }

        string num1 = exp[gen][spec - 1]
        , num2 = exp[gen][spec + 1];

        string cp = calc_params.dump();

        char* val = Subtract(strdup(&num1[0]), strdup(&num2[0]), strdup(&cp[0]), line);

        exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

        exp[gen].insert(exp[gen].begin() + spec - 1, val);
      }
    }

    while (expContain(exp, "=") || expContain(exp, "!=") || expContain(exp, "<") || expContain(exp, ">") || expContain(exp, "<=") || expContain(exp, ">=") || expContain(exp, "~~") || expContain(exp, "!~~") || expContain(exp, "~~~") || expContain(exp, "!~~~")) {
      vector<tuple<int, int>> indexes {expIndex(exp, "="), expIndex(exp, "!="), expIndex(exp, "<"), expIndex(exp, ">"),expIndex(exp, "<="), expIndex(exp, ">="), expIndex(exp, "~~"), expIndex(exp, "!~~"), expIndex(exp, "~~~"), expIndex(exp, "!~~~")};

      int min = 0;

      for (int i = 0; i < indexes.size(); i++) {
        int ngen, nspec, ogen, ospec;

        tie(ngen, nspec) = indexes[i];
        tie(ogen, ospec) = indexes[min];

        if (ngen == -1 || nspec == -1) continue;
        if (ogen || ospec) min = i;

        if (ogen > ngen) min = i;
        else if (ogen == ngen && ospec > nspec) min = i;
        else continue;
      }

      switch (min) {
        case 0: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "=");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (strcmp(ReturnInitC(&f[0]), ReturnInitC(&l[0])) == 0) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 1: {

            int gen, spec;

            tie(gen, spec) = expIndex(exp, "!=");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (!(strcmp(ReturnInitC(&f[0]), ReturnInitC(&l[0])) == 0)) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 2: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "<");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (IsLessC(&f[0], &l[0])) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 3: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, ">");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (!IsLessC(&f[0], &l[0]) && strcmp(ReturnInitC(&f[0]), ReturnInitC(&l[0])) != 0) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 4: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "<=");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (IsLessC(&f[0], &l[0]) || strcmp(ReturnInitC(&f[0]), ReturnInitC(&l[0])) == 0) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 5: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, ">=");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 1];

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 2);

            if (!IsLessC(&f[0], &l[0])) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");

          }
          break;
        case 6: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "~~");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 3]
            , dif = exp[gen][spec + 1];

            if (dif.rfind("-", 0) == 0) dif = dif.substr(1);

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 4);

            if (
              (
                IsLessC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line), &l[0])
                ||
                strcmp(ReturnInitC(&l[0]), ReturnInitC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line))) == 0
              )
              ||
              (
                IsLessC(Add(&f[0], &dif[0], &calc_params.dump()[0], line), &l[0])
                ||
                strcmp(ReturnInitC(&l[0]), ReturnInitC(Add(&f[0], &dif[0], &calc_params.dump()[0], line))) == 0
              )
            ) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");
          }
          break;
        case 7: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "!~~");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 3]
            , dif = exp[gen][spec + 1];

            if (dif.rfind("-", 0) == 0) dif = dif.substr(1);

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 4);

            if (!(
              (
                IsLessC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line), &l[0])
                ||
                strcmp(ReturnInitC(&l[0]), ReturnInitC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line))) == 0
              )
              ||
              (
                IsLessC(Add(&f[0], &dif[0], &calc_params.dump()[0], line), &l[0])
                ||
                strcmp(ReturnInitC(&l[0]), ReturnInitC(Add(&f[0], &dif[0], &calc_params.dump()[0], line))) == 0
              )
            )) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");
          }
          break;
        case 8: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "~~~");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 3]
            , dif = exp[gen][spec + 1];

            if (dif.rfind("-", 0) == 0) dif = dif.substr(1);

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 4);

            if (
              strcmp(ReturnInitC(Add(&f[0], &dif[0], &calc_params.dump()[0], line)), ReturnInitC(&l[0])) == 0
              ||
              strcmp(ReturnInitC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line)), ReturnInitC(&l[0])) == 0
            ) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");
          }
          break;
        case 9: {
            int gen, spec;

            tie(gen, spec) = expIndex(exp, "!~~~");

            string f = exp[gen][spec - 1]
            , l = exp[gen][spec + 3]
            , dif = exp[gen][spec + 1];

            if (dif.rfind("-", 0) == 0) dif = dif.substr(1);

            exp[gen].erase(exp[gen].begin() + spec - 1, exp[gen].begin() + spec + 4);

            if (
              IsLessC(Add(&f[0], &dif[0], &calc_params.dump()[0], line), &l[0])
              ||
              (
                IsLessC(&l[0], Subtract(&f[0], &dif[0], &calc_params.dump()[0], line))
                ||
                strcmp(ReturnInitC(&l[0]), ReturnInitC(Subtract(&f[0], &dif[0], &calc_params.dump()[0], line))) == 0
              )
            ) exp[gen].insert(exp[gen].begin() + spec - 1, "true");
            else exp[gen].insert(exp[gen].begin() + spec - 1, "false");
          }
          break;
      }

    }

    return exp;
  }
}

#endif
