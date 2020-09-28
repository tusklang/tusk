#ifndef SYSTABLES_SYSCALLS_STAT_H_
#define SYSTABLES_SYSCALLS_STAT_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
//will do windows later
#else
#include <sys/stat.h>

#define ulong unsigned long int

long int sysstat(char* file, ulong* st_dev, ulong* st_ino, ulong* st_mode) {
    struct stat* buf;
    buf->st_dev = *st_dev;
    buf->st_ino = *st_ino;
    buf->st_mode = *st_mode;

    //call the stat
    stat(file, buf);

    //pass back the pointers
    *st_dev = buf->st_dev;
    *st_ino = buf->st_ino;
    *st_mode = buf->st_mode;

}
#endif

#ifdef __cplusplus
}
#endif

#endif