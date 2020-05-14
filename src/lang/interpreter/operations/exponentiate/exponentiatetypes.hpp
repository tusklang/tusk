#ifndef EXPONENTIATETYPES_HPP_
#define EXPONENTIATETYPES_HPP_

#include <map>
#include <deque>
#include <vector>
#include "../../json.hpp"
#include "../../structs.hpp"
#include "../../values.hpp"
using namespace std;
using json = nlohmann::json;

Action exponentiatenumbers(Action num1, Action num2, json cli_params, deque<map<string, vector<Action>>> this_vals) {

  string n1 = num1.ExpStr[0], n2 = num2.ExpStr[0];

  string n2_(ReturnInitC(&n2[0]));

  //for now, if n2 is a float just return undef (CHANGE LATER to exp(a ln b))
  if (n2_.find(".") != string::npos) return falseyVal;

  bool neg = false;

  if (n2[0] == '-') {
    n2 = n2.substr(1);
    neg = true;
  }

  char* fin = "1";

  for (;((bool) IsLessC("0", &n2[0])); n2 = string(SubtractC(&n2[0], "1", &cli_params.dump()[0])))
    fin = MultiplyC(fin, &n1[0], &cli_params.dump()[0]);

  string finStr(fin);

  if (neg) finStr = "-" + finStr;

  return Action{ "number", "", { finStr }, emptyActVec, {}, emptyActVec2D, {}, 39, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private" };
}

#endif
