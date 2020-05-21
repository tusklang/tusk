#ifndef STRUCTS_HPP_
#define STRUCTS_HPP_

#include <vector>
#include <map>

namespace omm {

  //declare the structs
  struct Condition;
  struct Action;
  struct Variable;
  struct Returner;
  /////////////////////

  typedef struct Condition {

    std::string                                  Type;
    std::vector<Action>                          Condition;
    std::vector<Action>                          Actions;

  } Condition;

  typedef struct SubCaller {

    std::vector<std::vector<Action>>             Indexes;
    std::vector<std::vector<Action>>             Args;
    bool                                         IsProc;

  } SubCaller;

  typedef struct Action {

    std::string                                  Type;
    std::string                                  Name;
    std::vector<std::string>                     ExpStr;
    std::vector<Action>                          ExpAct;
    std::vector<std::string>                     Params;
    std::vector<std::vector<Action>>             Args;
    std::vector<Condition>                       Condition;

    int                                          ID;

    //stuff for operations

    std::vector<Action>                          First;
    std::vector<Action>                          Second;
    std::vector<Action>                          Degree;

    //stuff for indexes

    std::vector<std::vector<Action>>             Value;
    std::vector<std::vector<Action>>             Indexes;
    std::map<std::string, std::vector<Action>>   Hash_Values;

    bool                                         IsMutable;
    std::string                                  Access;
    std::vector<SubCaller>                       SubCall;

  } Action;

  typedef struct Variable {

    std::string                                  type;
    std::string                                  name;
    std::vector<Action>                          value;

  } Variable;

  typedef struct Returner {

    std::vector<std::string>                          value;
    std::map<std::string, Variable>              variables;
    Action                                       exp;
    std::string                                  type;

  } Returner;

}

#endif
