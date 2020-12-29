#ifndef SYSTABLES_SYSCALLS_PID_H_
#define SYSTABLES_SYSCALLS_PID_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <sys/types.h>
#ifdef _WIN32
#include <windows.h>
#define kill(pid, exitc) TerminateProcess((void *)pid, exitc)
#define gettid GetCurrentThreadId
#define tgkill(_, tid, exitc) TerminateThread((void *)tid, exitc)
#else
#include <unistd.h>
#include <sys/wait.h>
#include <signal.h>
#endif

    long long int sysgetpid()
    {
        return getpid();
    }

    long long int syswaitpid(long long int pid, long long int maxtime)
    {
#ifdef _WIN32
        return WaitForSingleObject(OpenProcess(PROCESS_ALL_ACCESS, TRUE, pid), maxtime);
#else
    return waitpid(pid, 0, 0);
#endif
    }

    long long int syskillpid(long long int pid, int exitc)
    {
        return kill(pid, exitc);
    }

    long long int sysgettid()
    {
        return gettid();
    }

    long long int systkill(long long int tid, int exitc)
    {
        return tgkill(-1, tid, exitc);
    }

#ifdef __cplusplus
}
#endif

#endif