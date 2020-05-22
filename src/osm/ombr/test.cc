#include <iostream>
#include "compile.hpp"
#include "../../lang/interpreter/values.hpp"

int main() {
  std::cout << omm::osm::ombr::compile(R"()", {}) << std::endl;
}
