#ifndef STRUCTS_HPP_
#define STRUCTS_HPP_

#include <vector>
#include <map>
using namespace std;

//declare the structs
struct Condition;
struct Action;
struct Variable;
struct Returner;
/////////////////////

typedef struct Condition {

  string                        Type;
  vector<Action>                Condition;
  vector<Action>                Actions;

} Condition;

typedef struct Action {

  string                        Type;
  string                        Name;
  vector<string>                ExpStr;
  vector<Action>                ExpAct;
  vector<string>                Params;
  vector<vector<Action>>        Args;
  vector<Condition>             Condition;

  int                           ID;

  //stuff for operations

  vector<Action>                First;
  vector<Action>                Second;
  vector<Action>                Degree;

  //stuff for indexes

  vector<vector<Action>>        Value;
  vector<vector<Action>>        Indexes;
  map<string, vector<Action>>   Hash_Values;

  bool                          IsMutable;

} Action;

typedef struct Variable {

  string                        type;
  string                        name;
  vector<Action>                value;

} Variable;

typedef struct Returner {

  vector<string>                value;
  map<string, Variable>         variables;
  Action                        exp;
  string                        type;

} Returner;

#endif
