#ifndef TUSK_NATIVE_SYSCALL_H_
#define TUSK_NATIVE_SYSCALL_H_

#ifdef __cplusplus
extern "C" {
#endif

const int MAX_SYSCALL_ARGC = 4;
#define tusksyscallargs \
    void* a0, \
    void* a1, \
    void* a2, \
    void* a3

#include "syscall_unix.h"
#include "syscall_win.h"

#ifdef __cplusplus
}
#endif

#endif