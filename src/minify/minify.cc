#include <string>
#include <fstream>
#include "replace.hpp"
using namespace std;

int main(int argc, char** argv) {

  string dir = argv[1];

  for (int i = 2; i < argc; i++) {
    ifstream fin;

    fin.open(dir + argv[i]);

    string file, line;

    while (fin) {
      getline(fin, line);

      file+=line;
    }

    string nFile = replaceWhitespace(file);

    ofstream fout;

    fout.open(dir + argv[i]);

    fout << nFile;
  }
}
