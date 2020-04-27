#ifndef SIMILARITY_HPP_
#define SIMILARITY_HPP_

#include "json.hpp"
using namespace std;

json similarity(json val1, json val2, json degree, const json calc_params, json vars, const string dir, int line) {

  //if the degree is not a number return undefined
  if (degree["Type"] != "number") return falseyVal;

  if (val1["Name"] == "hashed_value" && val2["Name"] == "hashed_values") {

    if (val1["Hash_Values"].size() < val2["Hash_Values"].size()) {
      json temp = val1;

      val1 = val2;
      val2 = temp;
    }

    char* difcount = "0";

    for (auto& i : val1["Hash_Values"].items()) {

      for (auto& o : val2["Hash_Values"].items()) {

        if (
          parser(o.value(), calc_params, vars, dir, false, line, true).exp["ExpStr"][0]
          !=
          parser(i.value(), calc_params, vars, dir, false, line, true).exp["ExpStr"][0]
        ) difcount = AddStrings(difcount, "1", &calc_params.dump()[0], line);
        else val2["Hash_Values"].erase(o.key());

        if ((bool) IsLessC(&(degree["ExpStr"][0].get<string>())[0], difcount)) return falseRet;
      }
    }

    //it has passed the above test
    return trueRet;

  } else {
    string
      num1E = val1["ExpStr"][0],
      num2E = val2["ExpStr"][0],
      degreeE = degree["ExpStr"][0];

      bool upperLess, lowerGreater;

      if (degreeE != "0") {
        char* upper = AddStrings(&num1E[0], &degreeE[0], &calc_params.dump()[0], line);
        char* lower = SubtractStrings(&num1E[0], &degreeE[0], &calc_params.dump()[0], line);

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

json strictSimilarity(json val1, json val2, json degree, const json calc_params, json vars, const string dir, int line) {

  //if the degree is not a number return undefined
  if (degree["Type"] != "number") return falseyVal;

  if (val1["Name"] == "hashed_value" && val2["Name"] == "hashed_values") {

    if (val1["Hash_Values"].size() < val2["Hash_Values"].size()) {
      json temp = val1;

      val1 = val2;
      val2 = temp;
    }

    char* difcount = "0";

    for (auto& i : val1["Hash_Values"].items()) {

      auto find = val2["Hash_Values"].find(i.key());

      if (find == val2["Hash_Values"].end()) difcount = AddStrings(difcount, "1", &calc_params.dump()[0], line);
      else {

        if (
          parser(*find, calc_params, vars, dir, false, line, true).exp["ExpStr"][0]
          !=
          parser(i.value(), calc_params, vars, dir, false, line, true).exp["ExpStr"][0]
        ) difcount = AddStrings(difcount, "1", &calc_params.dump()[0], line);
        else {
          val2["Hash_Values"].erase(i.key());
        }
      }

      if ((bool) IsLessC(&(degree["ExpStr"][0].get<string>())[0], difcount)) return falseRet;
    }

    //it has passed the above test
    return trueRet;

  } else {
    string
      num1E = val1["ExpStr"][0],
      num2E = val2["ExpStr"][0],
      degreeE = degree["ExpStr"][0];

      bool upperLess, lowerGreater;

      if (degreeE != "0") {
        char* upper = AddStrings(&num1E[0], &degreeE[0], &calc_params.dump()[0], line);
        char* lower = SubtractStrings(&degreeE[0], &num1E[0], &calc_params.dump()[0], line);

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
