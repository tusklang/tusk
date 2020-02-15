using namespace std;

int colorpick(string val) {
  if (val.rfind("'", 0) == 0 || val.rfind("\"", 0) == 0 || val.rfind("`", 0) == 0) return -1;
  else if (val == "undefined") return 8;
  
  return 6;
}
