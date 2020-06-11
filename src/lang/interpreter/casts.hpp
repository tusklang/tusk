#ifndef OMM_CASTER_HPP_
#define OMM_CASTER_HPP_

#include <map>
#include <functional>

#include "structs.hpp"
#include "values.hpp"
#include "operations/numeric/numeric.hpp"

namespace omm {

  //list of all of the convertions between types
  std::map<std::string, std::function<Action(Action val)>> casts = {
    { "number->string", [](Action val) -> Action { //number to string
      Action str = strPlaceholder;
      str.ExpStr[0] = normalize_number(val);
      return str;
    } }
  };

  Action cast(Action val, std::string type) {

    std::string castf = val.Type + "->" + type;

    if (casts.find(castf) == casts.end()) {

      //just switch the type in the action

      Action value = val;
      value.Type = type;
      return value;
    }

    return casts[castf](val);
  }

}

#endif
