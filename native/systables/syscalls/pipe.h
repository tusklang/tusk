#ifndef SYSTABLES_SYSCALLS_PIPE_H_
#define SYSTABLES_SYSCALLS_PIPE_H_

#ifdef __cplusplus
extern "C"
{
#endif

#ifdef _WIN32
#include <windows.h>
#endif
#include <unistd.h>

    long long int syspipe(void **fds, long long int size)
    {
#ifdef _WIN32
        return CreatePipe(&fds[0], &fds[1], NULL, size);
#else
    int nfds[2];
    nfds[0] = (long long int)fds[0];
    nfds[1] = (long long int)fds[1];
    return pipe(nfds);
#endif
    }

#ifdef __cplusplus
}
#endif

#endif