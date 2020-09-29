#ifndef SYSTABLES_SYSCALLS_FSTAT_H_
#define SYSTABLES_SYSCALLS_FSTAT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <sys/stat.h>

struct stat getstat(long int fd) {
    struct stat buf;
    fstat(fd, &buf);
    return buf;
}

//for DRY
#define getstatf(field) {return getstat(fd).field;}

long int fst_dev(long int fd) getstatf(st_dev)
long int fst_ino(long int fd) getstatf(st_ino)
long int fst_mode(long int fd) getstatf(st_mode)
long int fst_nlink(long int fd) getstatf(st_nlink)
long int fst_uid(long int fd) getstatf(st_uid)
long int fst_gid(long int fd) getstatf(st_gid)
long int fst_rdev(long int fd) getstatf(st_rdev)
long int fst_size(long int fd) getstatf(st_size)

#ifdef __cplusplus
}
#endif

#endif