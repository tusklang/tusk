#ifndef COMPARISONS_HPP_
#define COMPARISONS_HPP_

#include "json.hpp"
#include "values.hpp"
using namespace std;

json equals(json val1, json val2, json calc_params, int line) {

  if (val1["Name"] == "hashed_value" && val2["Name"] == "hashed_values") {

    if (val1["Hash_Values"].size() < val2["Hash_Values"].size()) {
      json temp = val1;

      val1 = val2;
      val2 = temp;
    }

    for (auto& i : val1["Hash_Values"].items()) {

      auto secondPart = val2["Hash_Values"].find(i.key());

      if (secondPart == val2["Hash_Values"].end()) return falseRet;
      if (*secondPart != i.value()) return falseRet;
    }

    return trueRet;
  } else {
    if (val1["ExpStr"][0] == val2["ExpStr"][0]) return trueRet;
    else return falseRet;
  }

  return falseRet;
}

#endif
