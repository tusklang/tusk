#ifndef STRUCTS_HPP_
#define STRUCTS_HPP_

#include <vector>
#include <map>
#include <deque>
#include <functional>
#include <future>
#include <memory>

namespace omm {

  //declare the structs
  struct Condition;
  struct Action;
  struct Action;
  struct Variable;
  struct Returner;
  /////////////////////

  typedef struct Condition {

    std::string                                     Type;
    std::vector<omm::Action>                        Condition;
    std::vector<omm::Action>                        Actions;

  } Condition;

  typedef struct SubCaller {

    std::vector<std::vector<omm::Action>>           Indexes;
    std::vector<std::vector<omm::Action>>           Args;
    bool                                            IsProc;

  } SubCaller;

  typedef struct Action {

    std::string                                     Type;
    std::string                                     Name;
    std::vector<std::string>                        ExpStr;
    std::vector<omm::Action>                        ExpAct;
    std::vector<std::string>                        Params;
    std::vector<std::vector<omm::Action>>           Args;
    std::vector<omm::Condition>                     Condition;

    int                                             ID;

    //stuff for operations

    std::vector<omm::Action>                        First;
    std::vector<omm::Action>                        Second;
    std::vector<omm::Action>                        Degree;

    //stuff for indexes

    std::vector<std::vector<omm::Action>>           Value;
    std::vector<std::vector<omm::Action>>           Indexes;
    std::map<std::string, std::vector<omm::Action>> Hash_Values;

    bool                                            IsMutable;
    std::string                                     Access;
    std::vector<omm::SubCaller>                     SubCall;

    std::vector<long long>                          Integer;
    std::vector<long long>                          Decimal;

    //values that are calculated at runtime
    std::shared_ptr<std::future<omm::Returner>>     Thread; //for async

  } Action;

  typedef struct Variable {

    std::string                                     type;
    std::string                                     name;
    std::vector<omm::Action>                        value;

    //cprocs
    std::function<Returner(
      omm::Action v,
      omm::CliParams cli_params,
      std::map<std::string, Variable> vars,
      std::deque<std::map<std::string, std::vector<Action>>> this_vals,
      std::string dir
    )>                                            cproc;

  } Variable;

  typedef struct Returner {

    std::vector<std::string>                      value;
    std::map<std::string, omm::Variable>          variables;
    omm::Action                                   exp;
    std::string                                   type;

  } Returner;

}

#endif
