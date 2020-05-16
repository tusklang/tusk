#ifndef DECODE_JSON_HPP_
#define DECODE_JSON_HPP_

#include <vector>
#include <map>
#include <functional>
#include "structs.hpp"
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

namespace DecodeJSON {

  Action action(json j);

  std::string string(json j) {
    return j.get<std::string>();
  }
  bool boolean(json j) {
    return j.get<bool>();
  }
  int integer(json j) {
    return j.get<int>();
  }
  std::vector<std::string> vectorStr(json j) {

    if (j.size() == 0 || j.is_null()) return {};

    return j.get<std::vector<std::string>>();
  }
  std::vector<Action> vector(json j) {

    std::vector<Action> ret;

    if (j.size() == 0 || j.is_null()) return {};

    for (json it : j) {

      Action val = action(it);
      ret.push_back(val);

    }

    return ret;
  }
  std::vector<std::vector<Action> > vector2D(json j) {

    std::vector<std::vector<Action>> ret;

    if (j.size() == 0 || j.is_null()) return {};

    for (json it : j) ret.push_back(vector(it));

    return ret;
  }
  std::vector<Condition> vectorCond(json j) {

    std::vector<Condition> ret;

    if (j.size() == 0 || j.is_null()) return {};

    for (json it : j) {

      Condition c;

      c.Type = string(it["Type"]);
      c.Condition = vector(it["Condition"]);
      c.Actions = vector(it["Actions"]);

      ret.push_back(c);

    }

    return ret;
  }
  std::vector<SubCaller> subcall(json j) {

    if (j.size() == 0 || j.is_null()) return {};

    std::vector<SubCaller> ret;

    for (json it : j) {

      SubCaller s;

      s.Indexes = vector2D(it["Indexes"]);
      s.Args = vector2D(it["Args"]);
      s.IsProc = boolean(it["IsProc"]);

      ret.push_back(s);

    }

    return ret;
  }
  std::map<std::string, std::vector<Action> > map(json j) {

    std::map<std::string, std::vector<Action>> ret;

    if (j.size() == 0 || j.is_null()) return {};

    for (auto& it : j.items()) ret[it.key()] = vector(it.value());

    return ret;
  }

  //all of functions (lambdas) to change the values of the action
  std::map<std::string, std::function<void(Action*, json)> > action_switch = {
    {"Type", [](Action* act, json j) { act->Type = string(j); } },
    {"Name", [](Action* act, json j) { act->Name = string(j); } },
    {"ExpStr", [](Action* act, json j) { act->ExpStr = vectorStr(j); } },
    {"ExpAct", [](Action* act, json j) { act->ExpAct = vector(j); } },
    {"Params", [](Action* act, json j) { act->Params = vectorStr(j); } },
    {"Args", [](Action* act, json j) { act->Args = vector2D(j); } },
    {"Condition", [](Action* act, json j) { act->Condition = vectorCond(j); } },
    {"ID", [](Action* act, json j) { act->ID = integer(j); } },
    {"First", [](Action* act, json j) { act->First = vector(j); } },
    {"Second", [](Action* act, json j) { act->Second = vector(j); } },
    {"Degree", [](Action* act, json j) { act->Degree = vector(j); } },
    {"Value", [](Action* act, json j) { act->Value = vector2D(j); } },
    {"Indexes", [](Action* act, json j) { act->Indexes = vector2D(j); } },
    {"Hash_Values", [](Action* act, json j) { act->Hash_Values = map(j); } },
    {"IsMutable", [](Action* act, json j) { act->IsMutable = boolean(j); } },
    {"Access", [](Action* act, json j) { act->Access = string(j); } },
    {"SubCall", [](Action* act, json j) { act->SubCall = subcall(j); } }
  };

  Action action(json j) {

    Action ret;

    for (auto& it : j.items())
      action_switch[it.key()](&ret, it.value());

    return ret;
  }

}

#endif
