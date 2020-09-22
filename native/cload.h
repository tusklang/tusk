#ifndef TUSK_NATIVE_CLOAD_H_
#define TUSK_NATIVE_CLOAD_H_

#ifdef __cplusplus
extern "C" {
#endif

//cload allows tusk to call c (c, c++, go, rust, etc...) libraries (dll/so)

const int MAX_ARGC = 20;

//include both, because the os is determined if both files
#include "cload_win.h"
#include "cload_unix.h"

#ifdef __cplusplus
}
#endif

#endif