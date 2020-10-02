#ifndef SYSTABLES_SYSCALLS_PID_H_
#define SYSTABLES_SYSCALLS_PID_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <windows.h>
#define kill(pid, exitc) TerminateProcess((void*) pid, exitc)
#else
#include <unistd.h>
#include <sys/wait.h>
#define WaitForSingleObject(fd, maxtime) wait((long int) fd);
#endif

long long int sysgetpid() {
    return getpid();
}

long long int syswaitpid(long long int pid, long long int maxtime) {
    return WaitForSingleObject((void*) pid, maxtime);
}

long long int syskillpid(long long int pid, int exitc) {
    return kill(pid, exitc);
}

#ifdef __cplusplus
}
#endif

#endif