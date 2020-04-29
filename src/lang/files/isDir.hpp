#ifndef ISDIR_HPP_
#define ISDIR_HPP_

#include <sys/types.h>
#include <sys/stat.h>
using namespace std;

bool isDir(string dir) {
  struct stat s;

  if (stat(&dir[0], &s) == 0) return s.st_mode & S_IFDIR;

  return false;
}

#endif
