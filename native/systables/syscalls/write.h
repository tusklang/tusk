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
    return write(fd, buf, size);
};

#ifdef __cplusplus
}
#endif

#endif