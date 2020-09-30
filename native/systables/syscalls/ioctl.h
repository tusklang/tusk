#ifndef SYSTABLES_SYSCALLS_IOCTL_H_
#define SYSTABLES_SYSCALLS_IOCTL_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <ioapiset.h>
#else
#include <sys/ioctl.h>
#endif

long long int sysioctl(long long int fd, long long int request, char* argp) {
    //no idea how to implment this
}

#ifdef __cplusplus
}
#endif

#endif