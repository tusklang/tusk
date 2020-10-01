#ifndef SYSTABLES_SYSCALLS_SHED_YIELD_H_
#define SYSTABLES_SYSCALLS_SHED_YIELD_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <sched.h>

#ifdef _WIN32
#include <processthreadsapi.h>
#define shed_yield SwitchToThread
#endif

long long int sysshedyield() {
    return shed_yield();
}

#ifdef __cplusplus
}
#endif

#endif