#ifndef SYSTABLES_SYSCALLS_OPEN_H_
#define SYSTABLES_SYSCALLS_OPEN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>

#ifdef _WIN32
#include <io.h>
#else
#include <fcntl.h>
#endif

long long int sysopen(char* name, char* mode) {
    return fileno(fopen(name, mode));
}

#ifdef __cplusplus
}
#endif

#endif