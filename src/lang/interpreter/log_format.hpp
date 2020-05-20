#ifndef LOG_FORMAT_HPP_
#define LOG_FORMAT_HPP_

#include <iostream>
#include <windows.h>
#include <stdio.h>
#include <vector>
#include <map>
#include "json.hpp"
#include "../bind.h"
using json = nlohmann::json;

namespace omm {

  void log_format(Action in, const json cli_params, std::map<std::string, Variable> vars, int hash_spacing, std::string doPrint) {

    if (in.Type == "hash") {
      std::map<std::string, vector<Action>> hashvals = in.Hash_Values;

      if (hashvals.size() == 0) cout << "[::]" << (doPrint == "print" ? "" : "\n");
      else {
        cout << "[:" << endl;

        for (std::pair<std::string, std::vector<Action>> it : hashvals) {
          std::string key = it.first;
          std::vector<Action> _value = it.second;

          cout << std::string(hash_spacing, ' ') << key << ": ";
          log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
        }

        cout << std::string(hash_spacing - 2, ' ') << ":]" << (doPrint == "print" ? "" : "\n");
      }
    } else if (in.Type == "array") {
      std::map<std::string, std::vector<Action>> hashvals = in.Hash_Values;

      if (hashvals.size() == 0) cout << "[]" << (doPrint == "print" ? "" : "\n");
      else {
        cout << "[" << endl;

        for (std::pair<std::string, std::vector<Action>> it : hashvals) {
          std::string key = it.first;
          std::vector<Action> _value = it.second;

          cout << std::string(hash_spacing, ' ') << key << ": ";
          log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
        }

        cout << std::string(hash_spacing - 2, ' ') << "]" << (doPrint == "print" ? "" : "\n");
      }
    } else if (in.Type == "process" || in.Type == "group") cout << "{PROCESS~ | GROUP~} " << "PARAM COUNT: " << in.Params.size() << (doPrint == "print" ? "" : "\n");
    else if (in.Name == "operation") {
      log_format(in.First[0], cli_params, vars, hash_spacing, "print");

      std::string op = in.Type;
      cout << " " << GetOp(&op[0]) << " ";
      log_format(in.Second[0], cli_params, vars, hash_spacing, "print");

    } else {

      string val = in.ExpStr[0];

      cout << val << (doPrint == "print" ? "" : "\n");
    }
  }

}

#endif
