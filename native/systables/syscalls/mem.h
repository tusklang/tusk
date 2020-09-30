#ifndef SYSTABLES_SYSCALLS_MEM_H_
#define SYSTABLES_SYSCALLS_MEM_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

long long int sysmalloc(long long int size) {
    return (long long int) malloc(size);
}

long long int sysfree(long long int ptr) {
    free((void*)(ptr));
    return 0;
}

#ifdef __cplusplus
}
#endif

#endif