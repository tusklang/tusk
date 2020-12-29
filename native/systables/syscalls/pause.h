#ifndef SYSTABLES_SYSCALLS_PAUSE_H_
#define SYSTABLES_SYSCALLS_PAUSE_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <unistd.h>
#ifdef _WIN32 /* windows does not have a pause() syscall */
#define pause() Sleep(INFINITE)
#endif

    long long int syspause()
    {
        pause();
        return 0;
    }

#ifdef __cplusplus
}
#endif

#endif