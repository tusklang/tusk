#include "interpreter/run.hpp"

void bindParser(char* actions, char* cli_params, char* dir, int argc, char ** argv) {
  run(actions, cli_params, dir, argc, argv);
}
