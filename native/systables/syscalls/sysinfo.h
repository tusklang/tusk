#ifndef SYSTABLES_SYSCALLS_SYSINFO_H_
#define SYSTABLES_SYSCALLS_SYSINFO_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <winsock.h>
#include <windows.h>
#else
#include <sys/utsname.h>
#endif
#include <stdlib.h>

long long int sysuname(char* sysname, char* nodename, char* release) {
    #ifdef _WIN32

    strcpy(sysname, "windows");
    gethostname(nodename, strlen(nodename));

    char* sysrel;

    int ver = GetVersion();
 
    int maj = (DWORD)(LOBYTE(LOWORD(ver)));
    int min = (DWORD)(HIBYTE(LOWORD(ver)));
    int build;

    if (ver < 0x80000000)              
        build = (DWORD)(HIWORD(ver));

    char *majs, *mins, *builds;
    itoa(maj, majs, 10);
    itoa(min, mins, 10);
    itoa(build, builds, 10);

    //put the data into sysrel
    int asize = strlen(majs) + strlen(mins) + strlen(builds) + 2; //calculate size to alloc
    sysrel = realloc(sysrel, asize); //alloc the required space
    sprintf(sysrel, "%d.%d.%d", maj, min, build);

    release = realloc(release, asize);
    strcpy(release, sysrel);
    free(sysrel);

    #else

    //for DRY
    #define allocsp(name) realloc(name, strlen(buf.name) * sizeof(char))
    #define sendback(name) name = strcpy(name, buf.name);
    /////////

    //get info in a buffer
    struct utsname* buf;
    uname(buf);

    //realloc the space
    allocsp(sysname);
    allocsp(nodename);
    allocsp(release);

    //send back all the data
    sendback(sysname);
    sendback(nodename);
    sendback(release);

    //free all the unused data now
    free(buf.sysname);
    free(buf.nodename);
    free(buf.release);
    #endif
}

#ifdef __cplusplus
}
#endif

#endif