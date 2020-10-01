#ifndef SYSTABLES_SYSCALLS_DUP_H_
#define SYSTABLES_SYSCALLS_DUP_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int sysdup(long long int fd) {
    return dup(fd);
}

long long int sysdup2(long long int fd, long long int nfd) {
    return dup2(fd, nfd);
}

#ifdef __cplusplus
}
#endif

#endif