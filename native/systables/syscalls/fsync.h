#ifndef SYSTABLES_SYSCALLS_FSYNC_H_
#define SYSTABLES_SYSCALLS_FSYNC_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>

#ifdef _WIN32
#include <fileapi.h>
#define sync _flushall
#endif

long long int syssync() {
    return sync();
}

long long int sysfsync(long int fd) {
    return fflush(fdopen(fd, "r+"));
}

#ifdef __cplusplus
}
#endif

#endif