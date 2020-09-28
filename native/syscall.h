#ifndef TUSK_NATIVE_SYSCALL_H_
#define TUSK_NATIVE_SYSCALL_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "systables/sysf.h"

static inline int makeintfromunsafe(void* v) {
    //prevent the warning, because it works
    #pragma GCC diagnostic ignored "-Wpointer-to-int-cast"
    #pragma GCC diagnostic push
    return (int) v;
    #pragma GCC diagnostic pop
}

static inline void* makeunsafeint(int v) {
    #pragma GCC diagnostic ignored "-Wint-to-pointer-cast"
    #pragma GCC diagnostic push
    return (void*) v;
    #pragma GCC diagnostic pop
}

static inline long int callsys(SYSF fn, void* a0, void* a1, void* a2, void* a3, void* a4, void* a5) {
    long int called = ((SYSF)(fn))(a0, a1, a2, a3, a4, a5); //call the sycall func
    return called; //return the val
}

#ifdef __cplusplus
}
#endif

#endif