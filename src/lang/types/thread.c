#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32

#include <stdlib.h>
#include <windows.h>
#include <stdbool.h>
#include "thread.h"

DWORD WINAPI threadf(LPVOID* tav) {
    struct ThreadArgs* ta = ((struct ThreadArgs*) tav);
    CallGoCB((*ta).gof, (*ta).output);
}

struct Thread newThread(unsigned long long ptr) {

    struct ThreadArgs ta;
    ta.gof = ptr;
    ta.output = (void**) malloc(1); //alloc to the heap

    struct ThreadArgs* taheap = calloc(1, sizeof(struct ThreadArgs));
    *taheap = ta;

    HANDLE h = CreateThread(NULL, 0, threadf, (LPVOID*) taheap, 0, NULL);

    struct Thread t;
    t.ta = taheap;
    t.handle = h;
    return t;
}

void freeInAndOut(struct Thread thread) {
    //used to deallocate the thread, when it goes out of scope
    //but not terminate the thread

    //free the storage
    free((*thread.ta).gof);
    free(*(*thread.ta).output);
    free((*thread.ta).output);
    free(thread.ta);
    //////////////////
}

bool closeThread(struct Thread thread, DWORD exitcode) {

    freeInAndOut(thread);

    return TerminateThread(thread.handle, exitcode);
}

void* waitfor(struct Thread t) {
    WaitForSingleObject(t.handle, INFINITE);
    return *(*t.ta).output;
}

DWORD getexitcode(struct Thread t) {
    bool ret = GetExitCodeThread(t.handle, &t.exitcode);

    if (!ret) return 1;
    return t.exitcode;
}

#else
#endif

#ifdef __cplusplus
}
#endif