#ifndef FALSEY_VAL_HPP_
#define FALSEY_VAL_HPP_

json falseyType = {
  {"Type", "falsey"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"undefined\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 41},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"ValueType", "[]"_json}
}
, falseyVal = {
  {"Type", "falsey"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"undefined\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 41},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"ValueType", json::parse("[" + falseyType.dump() + "]")}
};

#endif
