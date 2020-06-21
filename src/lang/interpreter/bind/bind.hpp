#ifndef OMM_BIND_BIND_HPP_
#define OMM_BIND_BIND_HPP_

//the max digit number (if digit surpasses this, overflow it)

/*

  The reason DigitSize must be 1 is because of zeros
  example:
    DigitSize = 2
    Number1 = 1000
    Number2 = 100

  1000 -> [10, 00] -> [10, 0]
  100 -> [10, 0] -> [10, 0]

  thus 1000 == 100, but that is not really true
*/

#define DigitSize 1 //if you change here, change in numconv.go

#define OMM_MAX_DIGIT (std::pow(10, DigitSize)) //the actual max digit + 1
#define OMM_MIN_DIGIT (-1 * OMM_MAX_DIGIT) //the actual min digit - 1

#include <iostream>
#include <vector>
#include <string>

extern "C" void Kill(); //kill the process
// extern "C" char* GetOp(char*); //get operation alias, e.g. add --> +, subtract --> -, etc.
// extern "C" char* ExecCmd(char*, char*, char*); //funtion to execute commands
// extern "C" std::vector<long long> DeNormalize(std::string); //function to "de-normalize" a number

#include "cli.hpp"
#include "../structs.hpp"
// #include "../parser.hpp"
// #include "../values.hpp"
// #include "../ommtypes.hpp"
// #include "../cprocs.hpp"

class BindClass {
  public:
    BindClass() {};
    void Run(std::vector<omm::Action>, omm::CliParams, std::string, std::vector<std::string>);
};

#endif
