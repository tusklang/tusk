#ifndef _WIN32

#ifndef TUSK_NATIVE_SYSCALL_UNIX_H_
#define TUSK_NATIVE_SYSCALL_UNIX_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <sys/syscall.h>
#include <errno.h>
#include <stdio.h>

static inline long int tusksyscall(void* a0, void* a1, void* a2, void* a3) {
    long int ret = syscall(
        (int) (*((double*)a0)), //eax is an int
        (int) (*((double*)a1)), //ebx is an int
        (int) (*((double*)a2)), //ecx is an int
        (*((char**)a3)) //edx is a char*
    );
    return ret;
}

#ifdef __cplusplus
}
#endif

#endif

#endif