#ifndef OSM_OMBR_CONNECTION_HPP_
#define OSM_OMBR_CONNECTION_HPP_

#include <vector>
#include <map>
#include <deque>

#include "../../lang/interpreter/parser.hpp"
#include "../../lang/interpreter/json.hpp"
#include "../../lang/interpreter/structs.hpp"
#include "../../lang/files/readfile.hpp"
#include "compile.hpp"
using json = nlohmann::json;

namespace omm {

  namespace osm {

    namespace ombr {

      //function to connect the omm language to osm::ombr
      char* connect(void* argsVoid, void* cli_paramsV, void* varsV, void* this_valsV, void* dirV) {

        std::vector<std::vector<Action>> args = *((std::vector<std::vector<Action>>*)argsVoid);

        //if the required amount of args is not met return "404 not found"
        if (args.size() != 2 && args.size() != 1) return "404 not found";

        std::string dir = (*((std::string*)(dirV)));

        //calculate the file directory
        std::string fileDir = parser(
          args[0],
          *((json*)(cli_paramsV)),
          *((std::map<std::string, Variable>*)(varsV)),
          false,
          true,
          *((std::deque<std::map<std::string, std::vector<Action>>>*)(this_valsV)),
          *((std::string*)(dirV))
        ).exp.ExpStr[0]
        , fileContents = readfile(
          &((dir.rfind("/") == (dir.size() - std::string("/").size()) ? dir.substr(0, dir.length() - 2) : dir) //if the dir ends with /, remove 1 from the end
          + "/public/" + fileDir)[0]); //read the file of the ombr file

        std::map<std::string, std::vector<Action>> replacerMap;

        if (args.size() == 2) //if there are 2 arguments
          //calculate the replacers
          replacerMap = parser(
            args[1],
            *((json*)(cli_paramsV)),
            *((std::map<std::string, Variable>*)(varsV)),
            false,
            true,
            *((std::deque<std::map<std::string, std::vector<Action>>>*)(this_valsV)),
            *((std::string*)(dirV))
          ).exp.Hash_Values;

        //mapped out to the strings
        std::map<std::string, std::string> replacerMapStr;

        //calculate the replacers (as strings)
        for (std::pair<std::string, std::vector<Action>> it : replacerMap)
          replacerMapStr[it.first] = parser( //parse the current iterator
            it.second,
            *((json*)(cli_paramsV)),
            *((std::map<std::string, Variable>*)(varsV)),
             false,
             true,
             *((std::deque<std::map<std::string, std::vector<Action>>>*)(this_valsV)),
             *((std::string*)(dirV))
           ).exp.ExpStr[0];

        //return the compiled (templated) ombr (json) file
        return &(compile(fileContents, replacerMapStr))[0];
      }

    }

  }

}

#endif
