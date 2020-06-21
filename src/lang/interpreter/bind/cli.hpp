#ifndef OMM_CLI_PARAMS_STRUCTS_HPP_
#define OMM_CLI_PARAMS_STRUCTS_HPP_

#include <string>

namespace omm {

  //structs for cli_params
  typedef struct Calc {
    int PREC;
    std::string O;
  } Calc;
  typedef struct Packages {
    std::string ADDON;
  } Package;
  typedef struct Files {
    std::string NAME;
    std::string DIR;
  } Files;
  class CliParams {
    public:
      Calc    Calc;
      Package Package;
      Files   Files;
  };
  ////////////////////////

}

#endif
