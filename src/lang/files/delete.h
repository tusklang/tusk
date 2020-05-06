#ifndef DELETE_H_
#define DELETE_H_

#ifdef __cplusplus
extern "C" {
#endif
#include <stdio.h>

void deletefile(char* dir) {

  //remove the file
  remove(dir);
}

#ifdef __cplusplus
}
#endif

#endif
