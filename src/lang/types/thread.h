#ifndef OMM_THREAD_H_
#define OMM_THREAD_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

#ifdef _WIN32

#include <windows.h>

struct ThreadArgs { //go function cb and output ptr
    unsigned long long  gof;
    void**              output;
};
struct Thread {
    struct ThreadArgs*  ta;
    DWORD               exitcode;
    HANDLE              handle;
};

#else

#endif

extern void CallGoCB(unsigned long long, void**);
struct Thread newThread(unsigned long long);
void freeInAndOut(struct Thread);
bool closeThread(struct Thread, DWORD);
void* waitfor(struct Thread);
DWORD getexitcode(struct Thread);

#ifdef __cplusplus
}
#endif

#endif