#ifndef SYSTABLES_SYSCALLS_SYNC_H_
#define SYSTABLES_SYSCALLS_SYNC_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <stdio.h>

#ifdef _WIN32
#include <fileapi.h>
#define sync flushall
#endif

long long int syssync() {
    return sync();
}

long long int syssyncfd(long int fd) {
    return fflush(fdopen(fd, "r+")); //flush the file
}

#ifdef __cplusplus
}
#endif

#endif