#ifndef SYSTABLES_SYSCALLS_ACCESS_H_
#define SYSTABLES_SYSCALLS_ACCESS_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <unistd.h>

    long long int sysaccess(char *path, int amode)
    {
        return access(path, amode);
    }

#ifdef __cplusplus
}
#endif

#endif