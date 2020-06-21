#include "bind.hpp"

void BindClass::Run(std::vector<omm::Action> acts, omm::CliParams params, std::string dir, std::vector<std::string> args) {

  Kill();

  std::cout << "I" << std::endl;

  // std::map<std::string, omm::Variable> vars;
  //
  // omm::Action dirnameAct = { "string", "", { dir }, emptyActVec, {}, emptyActVec2D, {}, 38, emptyActVec, emptyActVec, emptyActVec, emptyActVec2D, emptyActVec2D, noneMap, false, "private", emptySubCaller, emptyLLVec, emptyLLVec, emptyFuture };
  //
  // vars["$dirname"] = omm::Variable{
  //   "global",
  //   "$dirname",
  //   { dirnameAct },
  //   [](omm::Action v, omm::CliParams cli_params, std::map<std::string, omm::Variable> vars, std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals, std::string dir) -> omm::Returner { return omm::Returner{}; }
  // };
  //
  // omm::Variable omm_args = omm::Variable{
  //   "global",
  //   "$argv",
  //   { arrayVal },
  //   [](omm::Action v, omm::CliParams cli_params, std::map<std::string, omm::Variable> vars, std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals, std::string dir) -> omm::Returner { return omm::Returner{}; }
  // };
  //
  // for (int i = 0; i < args.size(); ++i)
  //   omm_args.value[0].Hash_Values[std::to_string(i)] = { ommtypes::to_string(args[i]) };
  //
  // vars["$argv"] = omm_args;
  //
  // //put all of the cprocs into the vars
  // for (std::pair<std::string, std::function<omm::Returner(
  //   omm::Action v,
  //   omm::CliParams cli_params,
  //   std::map<std::string, omm::Variable> vars,
  //   std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals,
  //   std::string dir
  // )>> it : cprocs)
  //   vars["$" + it.first] = omm::Variable{
  //     "cproc",
  //     "$" + it.first,
  //     {},
  //     it.second
  //   };
  //
  // parser(acts, params, vars, /*group return*/ false, /* expression return */ false, {}, dir);
}
