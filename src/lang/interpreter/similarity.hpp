#ifndef SIMILARITY_HPP_
#define SIMILARITY_HPP_

//header file to see if two values are similar

#include <map>
#include <vector>
#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
#include "structs.hpp"
using namespace std;
using json = nlohmann::json;

Action similarity(Action val1, Action val2, Action degree, const json cli_params, map<string, Variable> vars) {

  //if the degree is not a number return undef
  if (degree.Type != "number") return falseyVal;

  if (val1.Name == "hashed_value" && val2.Name == "hashed_values") {

    if (val1.Hash_Values.size() < val2.Hash_Values.size()) {

      //swap val1 and val2
      Action temp = val1;

      val1 = val2;
      val2 = temp;
    }

    char* difcount = "0";

    for (pair<string, vector<Action>> i : val1.Hash_Values) {

      for (pair<string, vector<Action>> o : val2.Hash_Values) {

        if (
          parser(o.second, cli_params, vars, false, true).exp.ExpStr[0]
          !=
          parser(i.second, cli_params, vars, false, true).exp.ExpStr[0]
        ) difcount = AddC(difcount, "1", &cli_params.dump()[0]);
        else val2.Hash_Values.erase(o.first);

        if ((bool) IsLessC(&(degree.ExpStr[0])[0], difcount)) return falseRet;
      }
    }

    //it has passed the above test
    return trueRet;

  } else {
    string
      num1E = val1.ExpStr[0],
      num2E = val2.ExpStr[0],
      degreeE = degree.ExpStr[0];

    bool upperLess, lowerGreater;

    if (degreeE != "0") {
      char* upper = AddC(&num1E[0], &degreeE[0], &cli_params.dump()[0]);
      char* lower = SubtractC(&num1E[0], &degreeE[0], &cli_params.dump()[0]);

      upperLess = ((bool) IsLessC(&num2E[0], upper)) || strcmp(ReturnInitC(&num2E[0]), ReturnInitC(upper)) == 0;
      lowerGreater = ((bool) IsLessC(lower, &num2E[0])) || strcmp(ReturnInitC(lower), ReturnInitC(&num2E[0])) == 0;
    } else {

      //if it is 0, no need to add (also it serves as lazy equality)
      upperLess = strcmp(ReturnInitC(&num2E[0]), ReturnInitC(&num1E[0])) == 0;
      lowerGreater = strcmp(ReturnInitC(&num1E[0]), ReturnInitC(&num2E[0])) == 0;
    }

    if (upperLess && lowerGreater) return trueRet;
    else return falseRet;
  }

  //if it reaches the end, return undefined
  return falseyVal;
}

Action strictSimilarity(Action val1, Action val2, Action degree, const json cli_params, map<string, Variable> vars) {

  //if the degree is not a number return undefined
  if (degree.Type != "number") return falseyVal;

  if (val1.Name == "hashed_value" && val2.Name == "hashed_values") {

    if (val1.Hash_Values.size() < val2.Hash_Values.size()) {
      Action temp = val1;

      val1 = val2;
      val2 = temp;
    }

    char* difcount = "0";

    for (pair<string, vector<Action>> i : val1.Hash_Values) {

      auto find = val2.Hash_Values.find(i.first);

      if (find == val2.Hash_Values.end()) difcount = AddC(difcount, "1", &cli_params.dump()[0]);
      else {

        if (
          parser(val2.Hash_Values[i.first], cli_params, vars, false, true).exp.ExpStr[0]
          !=
          parser(i.second, cli_params, vars, false, true).exp.ExpStr[0]
        ) difcount = AddC(difcount, "1", &cli_params.dump()[0]);
        else {
          val2.Hash_Values.erase(i.first);
        }
      }

      if ((bool) IsLessC(&degree.ExpStr[0][0], difcount)) return falseRet;
    }

    //it has passed the above test
    return trueRet;

  } else {
    string
      num1E = val1.ExpStr[0],
      num2E = val2.ExpStr[0],
      degreeE = degree.ExpStr[0];

      bool upperLess, lowerGreater;

      if (degreeE != "0") {
        char* upper = AddC(&num1E[0], &degreeE[0], &cli_params.dump()[0]);
        char* lower = SubtractC(&degreeE[0], &num1E[0], &cli_params.dump()[0]);

        //strict similarity for these values is just (+/-)
        upperLess = strcmp(ReturnInitC(&num2E[0]), ReturnInitC(upper)) == 0;
        lowerGreater = strcmp(ReturnInitC(lower), ReturnInitC(&num2E[0])) == 0;
      } else {

        //if it is 0, no need to add
        upperLess = strcmp(ReturnInitC(&num2E[0]), ReturnInitC(&num1E[0])) == 0;
        lowerGreater = strcmp(ReturnInitC(&num1E[0]), ReturnInitC(&num2E[0])) == 0;
      }

      if (upperLess || lowerGreater) return trueRet;
      else return falseRet;
  }

  //if it reaches the end, return undefined
  return falseyVal;
}

#endif
