#include <string>
#include <iostream>
using namespace std;

string replaceWhitespace(string file) {

  bool isEscaped = false;

  string nFile = ""
  , typeOfStr = "";

  for (int i = 0; i < file.length(); i++) {

    if (!isEscaped) {

      if (file[i] == '\\') {

        isEscaped = true;

        nFile+=file[i];
        continue;
      }

      if (file[i] == '\'') {
        if (typeOfStr == "") typeOfStr = "'";
        else if (typeOfStr == "'") typeOfStr = "";
      }

      if (file[i] == '\"') {
        if (typeOfStr == "") typeOfStr = "\"";
        else if (typeOfStr == "\"") typeOfStr = "";
      }

      if (file[i] == '\'') {
        if (typeOfStr == "") typeOfStr = "\'";
        else if (typeOfStr == "`") typeOfStr = "";
      }
    } else isEscaped = true;

    if (!(typeOfStr == "" && (file[i] == '\n' || file[i] == '\t' || file[i] == ' '))) {
      nFile+=file[i];
    } else if (typeOfStr == "" && !(file[i] == '\n' || file[i] == '\t' || file[i] == ' ')) {
      nFile+=file[i];
    }
  }

  cout << nFile;

  return nFile;
}
