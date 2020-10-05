#ifndef SYSTABLES_SYSCALLS_FILE_H_
#define SYSTABLES_SYSCALLS_FILE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

#ifdef _WIN32
#include <windows.h>
#include <winbase.h>
#define link(a, b) CreateHardLinkA(a, b, NULL)
#endif

long long int syslink(char* p1, char* p2) {
    return link(p1, p2);
}

long long int sysunlink(char* path) {
    return unlink(path);
}

long long int syschmod(char* name, int mode) {
    return chmod(name, mode);
}

#ifdef __cplusplus
}
#endif

#endif