#include <vector>
#include <map>
#include <deque>

#include "interpreter/parser.hpp"
#include "interpreter/structs.hpp"
#include "interpreter/json.hpp"
#include "interpreter/run.hpp"
using json = nlohmann::json;

void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv) {
  omm::run(actions, cli_params, dir, argc, argv);
}

void bindOsm(int handle_index, char* url) {

  std::deque<std::map<std::string, std::vector<omm::Action>>> this_vals;

  omm::Handler handle = omm::osm_handlers[handle_index];

  omm::Returner parsed = omm::parser(
    handle.callback.ExpAct,
    handle.cli_params,
    handle.vars,
    true,
    false,
    this_vals,
    handle.dir
  );
}
