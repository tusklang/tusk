#ifndef SYSTABLES_SYSCALLS_CLOSE_H_
#define SYSTABLES_SYSCALLS_CLOSE_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#include <winsock.h>
#else
#include <unistd.h>
#endif

long long int sysclose(long int fd) {
    int r = close(fd);

    #ifdef _WIN32
    //on windows, sockets cannot use the fs syscalls
    if (r == -1) return closesocket(fd);
    #endif

    return r;
}

#ifdef __cplusplus
}
#endif

#endif