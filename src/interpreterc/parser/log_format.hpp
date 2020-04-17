#ifndef LOG_FORMAT_HPP_
#define LOG_FORMAT_HPP_

#include <iostream>
#include "json.hpp"
#include "parser.hpp"
using namespace std;

Returner parser(const json actions, const json calc_params, json vars, const string dir, const bool groupReturn, int line, const bool expReturn);

void log_format(json in, const json calc_params, json vars, const string dir, int line, int hash_spacing, string doPrint) {

  if (in["Type"].dump() == "\"hash\"") {
    json hashvals = in["Hash_Values"];

    if (hashvals.size() == 0) cout << "[::]" << (doPrint == "print" ? "" : "\n");
    else {
      cout << "[:" << endl;

      for (json::iterator it = hashvals.begin(); it != hashvals.end(); it++) {
        json key = it.key()
        , _value = it.value();

        cout << string(hash_spacing, ' ') << key.dump().substr(1, key.dump().length() - 2) << ": ";
        log_format(_value[0], calc_params, vars, dir, line, hash_spacing + 2, "log");
      }

      cout << string(hash_spacing - 2, ' ') << ":]" << (doPrint == "print" ? "" : "\n");
    }
  } else if (in["Type"].dump() == "\"array\"") {
    json hashvals = in["Hash_Values"];

    if (hashvals.size() == 0) cout << "[]" << (doPrint == "print" ? "" : "\n");
    else {

      cout << "[" << endl;

      for (json::iterator it = hashvals.begin(); it != hashvals.end(); it++) {
        json key = it.key()
        , _value = it.value();

        cout << string(hash_spacing, ' ') << key.dump().substr(1, key.dump().length() - 2) << ": ";

        log_format(_value[0], calc_params, vars, dir, line, hash_spacing + 2, "log");
      }

      cout << "]" << (doPrint == "print" ? "" : "\n");
    }
  } else if (in["Type"].dump() == "\"process\"" || in["Type"].dump() == "\"group\"") cout << "{PROCESS~ | GROUP~}" << (doPrint == "print" ? "" : "\n");
  else if (in["Name"].dump() == "\"operation\"") {
    log_format(in["First"][0], calc_params, vars, dir, line, hash_spacing, "print");

    string op = in["Type"].get<string>();

    cout << " " << GetOp(&op[0]) << " ";
    log_format(in["Second"][0], calc_params, vars, dir, line, hash_spacing, "print");
  } else {

    string val = in["ExpStr"][0];

    cout << val << (doPrint == "print" ? "" : "\n");
  }
}

#endif
