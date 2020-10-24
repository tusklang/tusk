#ifndef SYSTABLES_SYSCALLS_WRITE_H_
#define SYSTABLES_SYSCALLS_WRITE_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#else
#include <unistd.h>
#endif

long long int syswrite(long int fd, char* buf, unsigned long long int size) {
    long int f = write(fd, buf, size);

    #ifdef _WIN32
    //because windows doesn't allow write for sockets
    if (f == -1) f = send(fd, buf, size, 0);
    #endif

    return f;
}

#ifdef __cplusplus
}
#endif

#endif