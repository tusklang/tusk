#ifndef SYSTABLES_SYSCALLS_OPEN_H_
#define SYSTABLES_SYSCALLS_OPEN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>
#include <dirent.h>
#include <sys/stat.h>

#ifdef _WIN32
#include <io.h>
#else
#include <fcntl.h>
#endif

long long int sysopen(char* name, char* mode) {

    //stat the file to see if it is a dir
    struct stat s;
    stat(name, &s);

    //if it is return the address of the DIR*
    if (!S_ISREG(s.st_mode)) return (long long int) opendir(name);
    //otherwise, just return the fd
    return fileno(fopen(name, mode));
}

#ifdef __cplusplus
}
#endif

#endif