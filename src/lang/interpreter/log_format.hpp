#ifndef LOG_FORMAT_HPP_
#define LOG_FORMAT_HPP_

#include <iostream>
#include <windows.h>
#include <stdio.h>
#include <vector>
#include <map>
#include "json.hpp"
#include "../bind.h"
using namespace std;
using json = nlohmann::json;

void log_format(Action in, const json cli_params, map<string, Variable> vars, int hash_spacing, string doPrint) {

  if (in.Type == "hash") {
    map<string, vector<Action>> hashvals = in.Hash_Values;

    if (hashvals.size() == 0) cout << "[::]" << (doPrint == "print" ? "" : "\n");
    else {
      cout << "[:" << endl;

      for (pair<string, vector<Action>> it : hashvals) {
        string key = it.first;
        vector<Action> _value = it.second;

        cout << string(hash_spacing, ' ') << key << ": ";
        log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
      }

      cout << string(hash_spacing - 2, ' ') << ":]" << (doPrint == "print" ? "" : "\n");
    }
  } else if (in.Type == "array") {
    map<string, vector<Action>> hashvals = in.Hash_Values;

    if (hashvals.size() == 0) cout << "[]" << (doPrint == "print" ? "" : "\n");
    else {
      cout << "[" << endl;

      for (pair<string, vector<Action>> it : hashvals) {
        string key = it.first;
        vector<Action> _value = it.second;

        cout << string(hash_spacing, ' ') << key << ": ";
        log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
      }

      cout << string(hash_spacing - 2, ' ') << "]" << (doPrint == "print" ? "" : "\n");
    }
  } else if (in.Type == "process" || in.Type == "group") cout << "{PROCESS~ | GROUP~} " << "PARAM COUNT: " << in.Params.size() << (doPrint == "print" ? "" : "\n");
  else if (in.Name == "operation") {
    log_format(in.First[0], cli_params, vars, hash_spacing, "print");

    string op = in.Type;
    cout << " " << GetOp(&op[0]) << " ";
    log_format(in.Second[0], cli_params, vars, hash_spacing, "print");

  } else {

    string val = in.ExpStr[0];

    cout << val << (doPrint == "print" ? "" : "\n");
  }
}

#endif
