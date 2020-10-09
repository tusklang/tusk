#ifndef SYSTABLES_SYSCALLS_READ_H_
#define SYSTABLES_SYSCALLS_READ_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#include <winsock.h>
#else
#include <unistd.h>
#endif

long long int sysread(long int fd, char* buf, unsigned long long int size) {
    int r = read(fd, buf, size);

    #ifdef _WIN32
    //on windows, sockets cannot use the fs syscalls
    if (r == -1) return recv(fd, buf, size, 0);
    #endif

    return r;
}

#ifdef __cplusplus
}
#endif

#endif