#ifndef SYSTABLES_SYSCALLS_SCHED_YIELD_H_
#define SYSTABLES_SYSCALLS_SCHED_YIELD_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <sched.h>

#ifdef _WIN32
#include <processthreadsapi.h>
#define sched_yield SwitchToThread
#endif

long long int sysschedyield() {
    return sched_yield();
}

#ifdef __cplusplus
}
#endif

#endif