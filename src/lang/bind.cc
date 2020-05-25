#include <vector>
#include <map>
#include <deque>
#include <thread>

#include "interpreter/parser.hpp"
#include "interpreter/structs.hpp"
#include "interpreter/json.hpp"
#include "interpreter/run.hpp"
#include "interpreter/ommtypes.hpp"
using json = nlohmann::json;

void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv) {
  omm::run(actions, cli_params, dir, argc, argv);
}

void bindOsm(int handle_index, char* url) {

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

  //make the osm handle func async
  std::thread t(
    omm::parser,
    handle.callback.ExpAct,
    handle.cli_params,
    vars,
    true,
    false,
    this_vals,
    handle.dir
  );
  t.detach();
}
