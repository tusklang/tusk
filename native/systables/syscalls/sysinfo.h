#ifndef SYSTABLES_SYSCALLS_SYSINFO_H_
#define SYSTABLES_SYSCALLS_SYSINFO_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <winsock.h>
#include <windows.h>
#include <sysinfoapi.h>
#include <psapi.h>
#else
#include <sys/sysinfo.h>
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
    void* _;
    #define allocsp(name) _ = realloc(name, strlen(buf->name) * sizeof(char))
    #define sendback(name) name = strcpy(name, buf->name);
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
    free(buf->sysname);
    free(buf->nodename);
    free(buf->release);
    #endif
}

long long int sysgetsysinfo(void** info) {

    /*
        0: time since last boot
        1: total ram
        2: free ram
        3: number of current processes
    */

    #ifdef _WIN32
    //windows

    info[0] = (void*) ((long long int) GetTickCount());

    //get ram stuff
    MEMORYSTATUSEX statex;
    statex.dwLength = sizeof(statex);
    GlobalMemoryStatusEx(&statex);
    ///////////////
    info[1] = (void*) statex.ullTotalPhys;
    info[2] = (void*) statex.ullAvailPhys;

    //get # of processes
    DWORD aProcesses[1024], cbNeeded;
    if (!EnumProcesses(aProcesses, sizeof(aProcesses), &cbNeeded)) return -1;
    info[3] = (void*) (cbNeeded / sizeof(DWORD));
    ////////////////////

    #elif defined TARGET_OS_X
    //mac

    //no idea how to do this
    //and I don't own a mac, so...

    #else
    //linux/bsd

    struct sysinfo* inf;
    int ret = sysinfo(inf);
    info[0] = (void*) inf->uptime;
    info[1] = (void*) inf->totalram;
    info[2] = (void*) inf->freeram;
    info[3] = (void*) ((long long int) inf->procs);

    return ret;
    #endif
}

#ifdef __cplusplus
}
#endif

#endif