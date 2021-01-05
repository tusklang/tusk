#ifndef TUSK_OPENLIB_H_
#define TUSK_OPENLIB_H_

#ifdef __cplusplus
extern "C"
{
#endif

    const int MAX_SYS_ARGC = 21; //max # of arguments to make a syscall or a native call

#include "openlib_win.h"
#include "openlib_unix.h"

#ifdef __cplusplus
}
#endif

#endif