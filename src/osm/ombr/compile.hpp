#ifndef OSM_OMBR_COMPILE_HPP_
#define OSM_OMBR_COMPILE_HPP_

#include <vector>
#include <map>
#include "../../lang/interpreter/json.hpp"
using json = nlohmann::json;

namespace omm {

  namespace osm {

    namespace ombr {

      std::string compile(std::string file_json, map<std::string, std::string> replacers) {

        if (file_json.length() == 0) return "";

        json j = json::parse(file_json);

        if (j.type() == json::value_t::object) j = json::parse("[" + file_json + "]");

        if (j.type() != json::value_t::array) return file_json; //if it is not an array, then return the file

        //actually replace the values
        for (auto& it : j.items()) {

          if (it.value().find("id") != it.value().end() && replacers.find(it.value()["id"]) != replacers.end())

            j[it.key()]["value"] = replacers["id"];

        }

        std::string html = "";

        //convert to html
        for (auto& it : j.items()) {

          //if it is not an object
          if (it.value().type() != json::value_t::object) continue;

          //if it is setting the doctype
          if (it.value()["tag"].get<std::string>() == "doctype") {
            html+="<!DOCTYPE " + it.value().get<std::string>();
            continue;
          }

          //create a tag e.g. <p id='example'>Ex</p>
          std::string tag = std::string("<") + it.value()["tag"].get<std::string>();

          //loop through the attrs
          for (auto& attrs : it.value().items()) {

            std::string compiled = compile(attrs.value().dump(), replacers);

            //add the attr to the tag
            tag+=std::string(" ") + attrs.key() + std::string("=") + compiled + std::string("");

          }

          //add the closing tag
          tag+=std::string("></") + it.value()["tag"].get<std::string>() + std::string("\">");

          //add the tag to the html
          html+=tag;

        }

        return html;
      }

    }

  }

}

#endif
