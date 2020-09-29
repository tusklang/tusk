#ifndef SYSTABLES_SYSCALLS_SYSPOLL_H_
#define SYSTABLES_SYSCALLS_SYSPOLL_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
//will do windows later
#else
#include <poll.h>

long int syspoll(long int fd, short events, short revents, unsigned long int nfds, int timeout) {
    struct pollfd pfd;
    pfd.fd = fd;
    pfd.events = events;
}

#endif

#ifdef __cplusplus
}
#endif

#endif