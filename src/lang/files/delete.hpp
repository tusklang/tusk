#ifndef DELETE_HPP_
#define DELETE_HPP_

#include "isDir.hpp"

#include <stdio.h>
#include <direct.h>
using namespace std;

void deletefile(string dir) {

  //convert to const char*
  const char* dirC = dir.c_str();

  if (isDir(dir))

    //remove the directory
    rmdir(dirC);
  else
  
    //remove the file
    remove(dirC);
}

#endif
