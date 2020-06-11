#ifndef OMM_SIMILARITY_HPP_
#define OMM_SIMILARITY_HPP_

//header file to see if two values are similar

#include <map>
#include <vector>
#include <deque>
#include "json.hpp"
#include "parser.hpp"
#include "values.hpp"
#include "structs.hpp"
#include "operations/numeric/numeric.hpp"
using json = nlohmann::json;

namespace omm {

  Action similarity(Action val1, Action val2, Action degree, const json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir);

  std::map<std::string, std::vector<Action>>::iterator findsimilar(std::map<std::string, std::vector<Action>> m, std::vector<Action> find, const json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    //loop through the map
    for (std::map<std::string, std::vector<Action>>::iterator it = m.begin(); it != m.end(); ++it) {

      if (
        similarity(
          parser(it->second, cli_params, vars, false, true, this_vals, dir).exp,
          parser(find, cli_params, vars, false, true, this_vals, dir).exp,
          zero,
          cli_params,
          vars,
          this_vals,
          dir
        ).ExpStr[0] == "true"
      ) return it;
    }

    return m.end();
  }

  Action similarity(Action val1, Action val2, Action degree, const json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    //if the degree is not a number return undef
    if (degree.Type != "number") return falseyVal;

    if (val1.Name == "hashed_value" && val2.Name == "hashed_value") {

      //force val1 to have a bigger hash than val2
      if (val1.Hash_Values.size() < val2.Hash_Values.size()) {
        Action temp = val1;

        val1 = val2;
        val2 = temp;
      }

      Action difcount = zero;

      std::map<std::string, std::vector<Action>> v1h = val1.Hash_Values, v2h = val2.Hash_Values;

      for (unsigned long long count = v1h.size(); count > 0; --count) {

        //get the first value from the first hash, and the same value from the other hash
        std::map<std::string, std::vector<Action>>::iterator
          v1find = v1h.begin(),
          v2find = findsimilar(v2h, v1find->second, cli_params, vars, this_vals, dir);

        if (v1find == v1h.end() || v2find == v2h.end())
          difcount = addNums(difcount, val1, &cli_params.dump()[0]);

        if (isLess(degree, difcount, cli_params)) return falseRet;
      }

      //it has passed the above test
      return trueRet;

    } else {

      bool upperLess, lowerGreater;

      if (!equals(degree, zero, cli_params)) {

        //get upper and lower bounds
        Action upper = addNums(val1, degree, cli_params);
        Action lower = subtractNums(val1, degree, cli_params);

        upperLess = isLess(val2, upper, cli_params) || equals(val2, upper, cli_params);
        lowerGreater = isLess(lower, val2, cli_params) || equals(lower, val2, cli_params);
      } else {

        //if it is 0, no need to add (also it serves as lazy equality)
        upperLess = equals(val2, val1, cli_params);
        lowerGreater = equals(val1, val2, cli_params);
      }

      if (upperLess && lowerGreater) return trueRet;
      else return falseRet;
    }

    //if it reaches the end, return undefined
    return falseyVal;
  }

  Action strictSimilarity(Action val1, Action val2, Action degree, const json cli_params, std::map<std::string, Variable> vars, std::deque<std::map<std::string, std::vector<Action>>> this_vals, std::string dir) {

    //if the degree is not a number return undefined
    if (degree.Type != "number") return falseyVal;

    if (val1.Name == "hashed_value" && val2.Name == "hashed_value") {

      //force val1 to have a bigger hash than val2
      if (val1.Hash_Values.size() < val2.Hash_Values.size()) {
        Action temp = val1;

        val1 = val2;
        val2 = temp;
      }

      Action difcount = zero;

      for (std::pair<std::string, std::vector<Action>> i : val1.Hash_Values) {

        auto find = val2.Hash_Values.find(i.first);

        if (find == val2.Hash_Values.end()) difcount = addNums(difcount, val1, &cli_params.dump()[0]);
        else {

          if (
            strictSimilarity(
              parser(val2.Hash_Values[i.first], cli_params, vars, false, true, this_vals, dir).exp,
              parser(i.second, cli_params, vars, false, true, this_vals, dir).exp,
              zero,
              cli_params,
              vars,
              this_vals,
              dir
            ).ExpStr[0] == "false"
          ) difcount = addNums(difcount, val1, &cli_params.dump()[0]);
          else {
            val2.Hash_Values.erase(i.first);
          }
        }

        if (isLess(degree, difcount, cli_params)) return falseRet;
      }

      //it has passed the above test
      return trueRet;

    } else {
      bool upperLess, lowerGreater;

      if (!equals(degree, zero, cli_params)) {
        Action upper = addNums(val1, degree, cli_params);
        Action lower = subtractNums(degree, val1, cli_params);

        //strict similarity for these values is just (+/-)
        upperLess = equals(val2, upper, cli_params);
        lowerGreater = equals(lower, val2, cli_params);
      } else {

        //if it is 0, no need to add
        upperLess = equals(val2, val1, cli_params);
        lowerGreater = equals(val1, val2, cli_params);
      }

      if (upperLess || lowerGreater) return trueRet;
      else return falseRet;
    }

    //if it reaches the end, return undefined
    return falseyVal;
  }

}

#endif
