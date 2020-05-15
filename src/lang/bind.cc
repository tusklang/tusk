#include "interpreter/run.hpp"

void bindParser(char* actions, char* cli_params, char* dir) {
  run(actions, cli_params, dir);
}
