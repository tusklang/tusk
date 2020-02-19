#include <string>
using namespace std;

bool endsWith(string s, string suffix) {
  return s.substr(s.length() - suffix.length()) == suffix;
}

string format(string val) {
  if (val.rfind("'", 0) == 0 || val.rfind("\"", 0) == 0 || val.rfind("`", 0) == 0) val = val.substr(1);
  if (endsWith(val, "'") || endsWith(val, "\"") || endsWith(val, "`")) val = val.substr(0, val.length() - 1);

  return val;
}
