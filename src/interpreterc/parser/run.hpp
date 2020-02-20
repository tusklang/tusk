#include "json.hpp"
#include "parser.hpp"
using namespace std;
using json = nlohmann::json;

void run(char *acts, char *calc_params, char *dir) {
  json actions = acts;
  json cp = calc_params;

  parser(actions, cp, (string) dir);
}
