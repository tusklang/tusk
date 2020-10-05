#ifndef SYSTABLES_SYSCALLS_LSTAT_H_
#define SYSTABLES_SYSCALLS_LSTAT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int syslseek(long int fd, long int offset, int whence) {
    return lseek(fd, offset, whence);
}

#ifdef __cplusplus
}
#endif

#endif