#ifndef SYSTABLES_SYSCALLS_FSTAT_H_
#define SYSTABLES_SYSCALLS_FSTAT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <sys/stat.h>

long int sysfsize(long int fd) {
    struct stat buf;
    fstat(fd, &buf);
    return buf.st_size; //return the size of buf
}

#ifdef __cplusplus
}
#endif

#endif