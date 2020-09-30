#ifndef SYSTABLES_SYSCALLS_MUNMAP_H_
#define SYSTABLES_SYSCALLS_MUNMAP_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <memoryapi.h>
#define munmap(addr, length) UnmapViewOfFile(addr);
#else
#include <unistd.h>
#include <sys/mman.h>
#endif

long long int sysmunmap(long long int addr, long long int length) {
    return munmap((void*)(addr), length);
}

#ifdef __cplusplus
}
#endif

#endif