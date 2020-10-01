#ifndef SYSTABLES_SYSCALLS_PID_H_
#define SYSTABLES_SYSCALLS_PID_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int sysgetpid() {
    return getpid();
}

#ifdef __cplusplus
}
#endif

#endif