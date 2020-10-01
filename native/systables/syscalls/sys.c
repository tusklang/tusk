#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32 //because winsock2 must be included before windows.h
#include <winsock2.h>
#include <windows.h>
#endif

//include all of the syscalls
#include "read.h"
#include "write.h"
#include "open.h"
#include "close.h"
#include "fstat.h"
#include "lseek.h"
#include "ioctl.h"
#include "read_writev.h"
#include "pipe.h"
#include "mem.h"
#include "select.h"
#include "sched_yield.h"
#include "dup.h"
#include "pause.h"
/////////////////////////////

#ifdef __cplusplus
}
#endif