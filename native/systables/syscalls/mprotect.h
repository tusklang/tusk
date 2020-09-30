#ifndef SYSTABLES_SYSCALLS_MPROTECT_H_
#define SYSTABLES_SYSCALLS_MPROTECT_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <memoryapi.h>
#else
#include <sys/mman.h>
#define VirtualProtect(addr, dwSize, flNewProtect, lpflOldProtect) mprotect(addr, dwSize, flNewProtect)
#define DWORD unsigned long int
#endif

long long int sysmprotect(long long int addr, int dwSize, long long int flNewProtect, long int lpflOldProtect) {
    DWORD oldprotv = ((DWORD) lpflOldProtect);
    return VirtualProtect((void*) addr, dwSize, flNewProtect, &oldprotv);
}

#ifdef __cplusplus
}
#endif

#endif