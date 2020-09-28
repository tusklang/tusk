#ifndef SYSTABLES_SYSCALLS_READ_H_
#define SYSTABLES_SYSCALLS_READ_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <io.h>
#else
#include <unistd.h>
#endif

long int sysread(long int fd, char* buf, unsigned long long int size) {
    return read(fd, buf, size);
};

#ifdef __cplusplus
}
#endif

#endif