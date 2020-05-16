#ifndef READFILE_HPP_
#define READFILE_HPP_

#include <vector>
#include <fstream>
#include <sstream>
#include <iostream>
#include <direct.h>
using namespace std;

string readfile(char* dir) {

  ifstream f(dir);

  stringstream ss;
  ss << f.rdbuf();

  string file = ss.str();

  return file;
}

#endif
