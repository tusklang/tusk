#include "interpreter/run.hpp"
#include "bind.h"
#include <iostream>
using namespace std;

void bindCgo(char *actions, char *cli_params) {
  run(actions, cli_params);
}
