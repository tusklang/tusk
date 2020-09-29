#ifndef SYSTABLES_SYSCALLS_MMAP_H_
#define SYSTABLES_SYSCALLS_MMAP_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <memoryapi.h>
#define mmap(addr, length, prot, flags, fd, offset) MapViewOfFile(addr, fd, 0, offset, length);
#else
#include <unistd.h>
#include <sys/mman.h>
#endif

long long int sysmmap(long long int addr, unsigned long long int length, int prot, int flags, int fd, long int offset) {
    void* p = mmap((void*) addr, length, prot, flags, fd, offset);
    return (long long int) p; //return the address of p
}

#ifdef __cplusplus
}
#endif

#endif