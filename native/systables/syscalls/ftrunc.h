#ifndef SYSTABLES_SYSCALLS_FTRUNC_H_
#define SYSTABLES_SYSCALLS_FTRUNC_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int sysftrucate(long int fd, long long int length) {
    return ftruncate(fd, length);
}

#ifdef __cplusplus
}
#endif

#endif