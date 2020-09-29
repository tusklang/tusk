#ifndef SYSTABLES_SYSCALLS_STAT_H_
#define SYSTABLES_SYSCALLS_STAT_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
//will do windows later
#else
#include <sys/stat.h>

long int sysstat(char* file) {
    struct stat* buf;
    stat(file, buf);
    
    //prevent the warning that gets the mem addr
    #pragma GCC diagnostic ignored "-Wpointer-to-int-cast"
    #pragma GCC diagnostic push
    return (int) buf; //return the mem addr of buf
    #pragma GCC diagnostic pop
}

#endif

#ifdef __cplusplus
}
#endif

#endif