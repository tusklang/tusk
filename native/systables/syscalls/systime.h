#ifndef SYSTABLES_SYSCALLS_SYSTIME_H_
#define SYSTABLES_SYSCALLS_SYSTIME_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>
#include <time.h>
#include <sys/time.h>

#ifdef _WIN32
#include <windows.h>
#endif

long long int sysgettime() {
    return time(NULL);
}

long long int sysgettimezone() {
    tzset();
    return timezone;
}

long long int syssettime(long long int unixtime) {
    #ifdef _WIN32

    //slightly modified from here:
    //  https://stackoverflow.com/questions/11122647/how-can-i-create-a-systemtime-struct-from-ulonglong-milliseconds/11123106#11123106

    time_t multiplier = 10000;
    time_t t = multiplier * unixtime;

    ULARGE_INTEGER li;
    li.QuadPart = t;

    FILETIME ft;
    ft.dwLowDateTime = li.LowPart;
    ft.dwHighDateTime = li.HighPart;

    SYSTEMTIME* syst;

    FileTimeToSystemTime(&ft, syst);
    return SetSystemTime(syst);
    #else
    struct timeval t;
    t.tv_usec = 1000 * unixtime;
    return settimeofday(&t, NULL);
    #endif
}

long long int syssettimezone(long long int lgmt) {
    tzset();
    timezone = lgmt;
    return 0;
}

#ifdef __cplusplus
}
#endif

#endif