#include <iostream>
#include "printer.hpp"
#include "colorpick.hpp"
#include "printformat.hpp"
using namespace std;

int main() {

  string val = "";

  for (string line; getline(cin, line);) val+=line;

  int color = colorpick(val);
  string formatted = format(val);
  printer(color, formatted);

  return 0;
}
