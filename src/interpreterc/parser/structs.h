#ifndef STRUCTS_H_
#define STRUCTS_H_

#ifdef __cplusplus
extern "C" {
#endif
  #include <vector>
  #include "json.hpp"
  using namespace std;
  using json = nlohmann::json;

  struct Variable {
    string                  type;
    string                  name;
    vector<string>          value;
    json                    valueActs;
  };

  struct Returner {
    vector<string>          value;
    json                    variables;
    json                    exp;
    string                  type;
  };
#ifdef __cplusplus
}
#endif

#endif
