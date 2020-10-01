#ifndef SYSTABLES_SYSCALLS_MMAP_H_
#define SYSTABLES_SYSCALLS_MMAP_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <memoryapi.h>
#define mmap(addr, length, prot, flags, fd, offset) MapViewOfFile(addr, fd, 0, offset, length);
#define munmap(addr, length) UnmapViewOfFile(addr);
#else
#include <unistd.h>
#include <sys/mman.h>
#define DWORD unsigned long int
#define VirtualProtect(addr, dwSize, flNewProtect, lpflOldProtect) mprotect(addr, dwSize, flNewProtect)
#endif

long long int sysmmap(long long int addr, unsigned long long int length, int prot, int flags, int fd, long int offset) {
    void* p = mmap((void*) addr, length, prot, flags, fd, offset);
    return (long long int) p; //return the address of p
}

long long int sysmprotect(long long int addr, int dwSize, long long int flNewProtect, long int lpflOldProtect) {
    DWORD oldprotv = ((DWORD) lpflOldProtect);
    return VirtualProtect((void*) addr, dwSize, flNewProtect, &oldprotv);
}

long long int sysmunmap(long long int addr, long long int length) {
    return munmap((void*)(addr), length);
}

#ifdef __cplusplus
}
#endif

#endif