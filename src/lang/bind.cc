#include "interpreter/run.hpp"
using json = nlohmann::json;

//function to bind c++ and go for the interpreter
void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv) {
  omm::run(actions, cli_params, dir, argc, argv);
}
