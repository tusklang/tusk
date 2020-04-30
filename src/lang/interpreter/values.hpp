#ifndef FALSEY_VAL_HPP_
#define FALSEY_VAL_HPP_

#include "json.hpp"
using namespace std;

const json falseyVal = {
  {"Type", "falsey"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"undef\"]")},
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
  {"IsMutable", false}
}
, trueRet = {
  {"Type", "boolean"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"true\"]")},
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
  {"IsMutable", false}
}
, falseRet = {
  {"Type", "boolean"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"false\"]")},
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
  {"IsMutable", false}
}
, zero = {
  {"Type", "number"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"0\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 39},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"IsMutable", false}
}
, val1 = {
  {"Type", "number"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"1\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 39},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"IsMutable", false}
}
, valn1 = {
  {"Type", "number"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"-1\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 39},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"IsMutable", false}
}
, strPlaceholder = {
  {"Type", "string"},
  {"Name", ""},
  {"ExpStr", json::parse("[\"\"]")},
  {"ExpAct", "[]"_json},
  {"Params", "[]"_json},
  {"Args", "[]"_json},
  {"Condition", "[]"_json},
  {"ID", 38},
  {"First", "[]"_json},
  {"Second", "[]"_json},
  {"Degree", "[]"_json},
  {"Value", "[[]]"_json},
  {"Indexes", "[[]]"_json},
  {"Index_Type", ""},
  {"Hash_Values", "{}"_json},
  {"IsMutable", false}
};

#endif
