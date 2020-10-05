#ifndef SYSTABLES_SYSCALLS_ACCESS_H_
#define SYSTABLES_SYSCALLS_ACCESS_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <time.h>

#ifdef _WIN32
#include <windows.h>
#define sleep Sleep
#endif

long long int syssleep(long long int amt) {
    sleep(amt);
    return 0;
}

#ifdef __cplusplus
}
#endif

#endif