#include "bind.h"
#include "parser/run.hpp"
#include <iostream>
using namespace std;

void bind(char *actions, char *calc_params, char *dir) {
  run(actions, calc_params, dir);
}

int main() {

  string data = "[{\"Type\":\"newline\",\"Name\":\"\",\"ExpStr\":[],\"ExpAct\":[],\"Params\":[],\"Args\":[],\"Condition\":[],\"Indexes\":[],\"ID\":0},{\"Type\":\"newline\",\"Name\":\"\",\"ExpStr\":[],\"ExpAct\":[],\"Params\":[],\"Args\":[],\"Condition\":[],\"Indexes\":[],\"ID\":0},{\"Type\":\"expression\",\"Name\":\"\",\"ExpStr\":[\"(\",\"true\",\")\"],\"ExpAct\":[],\"Params\":[],\"Args\":[],\"Condition\":[],\"Indexes\":[],\"ID\":7},{\"Type\":\"newline\",\"Name\":\"\",\"ExpStr\":[],\"ExpAct\":[],\"Params\":[],\"Args\":[],\"Condition\":[],\"Indexes\":[],\"ID\":0}]";

  char *d = &data[0];

  run(d, "{}", "./");
}
