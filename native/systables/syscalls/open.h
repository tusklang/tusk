#ifndef SYSTABLES_SYSCALLS_OPEN_H_
#define SYSTABLES_SYSCALLS_OPEN_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#else
#include <fcntl.h>
#endif

long int sysopen(char* name, int mode) {
    return open(name, mode);
};

#ifdef __cplusplus
}
#endif

#endif