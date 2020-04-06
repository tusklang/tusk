#include "parser/run.hpp"
#include "bind.h"
#include <iostream>
using namespace std;

void bind(char *actions, char *calc_params, char *dir) {
  run(actions, calc_params, dir);
}
