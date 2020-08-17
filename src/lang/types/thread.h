#ifndef OMM_THREAD_H_
#define OMM_THREAD_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

#ifdef _WIN32

#include <windows.h>

struct ThreadArgs {
    void* gof;
    void* output;
};
struct Thread {
    struct ThreadArgs ta;
    DWORD      exitcode;
    HANDLE     handle;
};

#else

#endif

extern void* Callgointerpreter(void*);
struct Thread newThread(void* cb);
void freeInAndOut(struct Thread thread);
bool closeThread(struct Thread thread, DWORD exitcode);
void* waitfor(struct Thread t);
DWORD getexitcode(struct Thread t);

#ifdef __cplusplus
}
#endif

#endif