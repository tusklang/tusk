#ifndef SYSTABLES_SYSCALLS_EXECV_H_
#define SYSTABLES_SYSCALLS_EXECV_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

long long int sysexecv(char* path, void** argv, void** newenviron) {
    return execv(path, (char**) argv 
        #ifndef _WIN32
        ,(char**) newenviron
        #endif
    );
}

#ifdef __cplusplus
}
#endif

#endif