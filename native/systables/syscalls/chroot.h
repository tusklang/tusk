#ifndef SYSTABLES_SYSCALLS_CHROOT_H_
#define SYSTABLES_SYSCALLS_CHROOT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int syschroot(char* path) {
    #ifndef _WIN32
    //not to be used on windows!!
    return chroot(path);
    #endif
    return -1;
}

#ifdef __cplusplus
}
#endif

#endif