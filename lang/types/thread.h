#ifndef OMM_THREAD_H_
#define OMM_THREAD_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

struct ThreadArgs { //go function cb and output ptr
    unsigned long long  gof;
    void**              output;
};

#ifdef _WIN32

#include <windows.h>

struct Thread {
    struct ThreadArgs*  ta;
    DWORD               exitcode;
    HANDLE              handle;
};

#else

#include <pthread.h>

struct Thread {
    struct ThreadArgs*  ta;
    unsigned long       exitcode;
    pthread_t           handle;
};

#endif

extern void CallGoCB(unsigned long long, void**);
struct Thread newThread(unsigned long long);
void freeThread(struct Thread);
void* waitfor(struct Thread);

#ifdef __cplusplus
}
#endif

#endif