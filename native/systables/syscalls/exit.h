#ifndef SYSTABLES_SYSCALLS_EXIT_H_
#define SYSTABLES_SYSCALLS_EXIT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int sysexit(int code) {
    exit(code);
}

#ifdef __cplusplus
}
#endif

#endif