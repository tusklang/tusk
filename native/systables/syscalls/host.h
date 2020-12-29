#ifndef SYSTABLES_SYSCALLS_HOST_H_
#define SYSTABLES_SYSCALLS_HOST_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <unistd.h>

#ifdef _WIN32
#include <winsock.h>
//I dont know if this is the same function, please fix if not
#define sethostname(name, len) SetComputerNameA(name)
#endif

    long long int sysgethostname(char *name, long long int len)
    {
        return gethostname(name, len);
    }

    long long int syssethostname(char *name, long long int len)
    {
        return sethostname(name, len);
    }

    long long int sysgetdomainname(char *name, long long int len)
    {
#ifndef _WIN32
        //not on windows
        return getdomainname(name, len);
#endif
        return -1;
    }

    long long int syssetdomainname(char *name, long long int len)
    {
#ifndef _WIN32
        //not on windows
        return setdomainname(name, len);
#endif
        return -1;
    }

#ifdef __cplusplus
}
#endif

#endif