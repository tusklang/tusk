#ifndef ISFILE_HPP_
#define ISFILE_HPP_

#include <sys/types.h>
#include <sys/stat.h>
using namespace std;

bool isFile(string dir) {
  struct stat s;

  if (stat(&dir[0], &s) == 0) return s.st_mode & S_IFREG;

  return false;
}

#endif
