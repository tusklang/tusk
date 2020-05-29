#include <vector>
#include <map>
#include <deque>
#include <thread>

#include "interpreter/parser.hpp"
#include "interpreter/structs.hpp"
#include "interpreter/json.hpp"
#include "interpreter/run.hpp"
#include "interpreter/ommtypes.hpp"
#include "../osm/ombr/connect.hpp"
#include "../osm/osm_render_alloc.h"
using json = nlohmann::json;

//function to bind c++ and go for the interpreter
void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv) {
  omm::run(actions, cli_params, dir, argc, argv);
}

//function to bind to c++ with go for osm
void bindOsm(int handle_index, char* url, osmGoProc goprocs[], osmGoProcName goprocNames[], int goprocsLen /*length of the goprocs*/) {

  std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals;

  omm::Handler handle = omm::osm_handlers[handle_index];

  //set osm variables
  std::map<std::string, omm::Variable> vars = handle.vars;

  //send the request url
  vars["$req.url"] = {
    "local",
    "$req.url",
    { omm::ommtypes::to_string(std::string(url)) },
    [](omm::Action v, json cli_params, std::map<std::string, omm::Variable> vars, std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals, std::string dir) -> omm::Returner { return omm::Returner{}; }
  };

  //account for all of the goprocs, e.g. render(), cookie.get(), cookie.set()
  for (int i = 0; i < goprocsLen; ++i) {
    vars["$" + std::string(goprocNames[i])] = omm::Variable{
      "cproc",
      "$" + std::string(goprocNames[i]),
      {},
      [i, goprocs](omm::Action v, json cli_params, std::map<std::string, omm::Variable> vars, std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals, std::string dir) -> omm::Returner {

        std::vector<std::string> retNo;

        return omm::Returner{ retNo, vars, omm::ommtypes::to_string(std::string(CallOSMProc(goprocs[i], ((void*) &v.Args), ((void*) &cli_params), ((void*) &vars), ((void*) &this_vals), ((void*) &dir)))), "expression" };
      }
    };
  }

  //make the osm handle func async
  //by making a new thread to call the osm code
  omm::parser(
    handle.callback.ExpAct,
    handle.cli_params,
    vars,
    true,
    false,
    this_vals,
    handle.dir
  );
  // t.detach();
}

char ** parseParams(void* params, void* cli_params, void* vars, void* this_vals, void* dir, int amt) {

  //convert args into a vector of a vector of actions
  std::vector<std::vector<omm::Action>> args = *(((std::vector<std::vector<omm::Action>>*)(params)));

  char ** parsedArgs = (char **) malloc(amt);
  int index = 0;

  if (amt != args.size()) return parsedArgs;

  for (std::vector<omm::Action> it : args) {
    parsedArgs[index] = &parser(
      args[0],
      *((json*)(cli_params)),
      *((std::map<std::string, omm::Variable>*)(vars)),
      false,
      true,
      *((std::deque<std::map<std::string, std::vector<omm::Action>>>*)(this_vals)),
      *((std::string*)(dir))
    ).exp.ExpStr[0][0]; //parse the current iterator

    ++index;
  }

  return parsedArgs;
}

//function to bind ombr
char* ombrBind(void* args, void* cli_params, void* vars, void* this_vals, void* dir) {
  return omm::osm::ombr::connect(args, cli_params, vars, this_vals, dir);
}
