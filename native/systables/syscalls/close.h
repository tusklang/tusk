#ifndef SYSTABLES_SYSCALLS_CLOSE_H_
#define SYSTABLES_SYSCALLS_CLOSE_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#else
#include <unistd.h>
#endif

long long int sysclose(long int fd) {
    return close(fd);
};

#ifdef __cplusplus
}
#endif

#endif