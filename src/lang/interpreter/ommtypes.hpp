#ifndef OMMTYPES_HPP_
#define OMMTYPES_HPP_

namespace omm {

  namespace ommtypes {

    #include "structs.hpp"
    #include "values.hpp"

    //func to convert std::string to omm string
    Action to_string(std::string str) {

      Action omm_str = strPlaceholder;

      omm_str.ExpStr[0] = str;

      //loop through the std::string to get hash values
      unsigned long long i = 0;
      for (char c : str) {

        Action curChar = strPlaceholder;
        curChar.ExpStr[0] = std::string(c, 1);

        omm_str.Hash_Values[std::to_string(i)] = { curChar };

        ++i;
      }

      return omm_str;
    }

  }

}

#endif
