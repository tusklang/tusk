#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32

#include <stdlib.h>
#include <windows.h>
#include <stdbool.h>
#include "thread.h"

DWORD WINAPI threadf(void* t) {
    struct ThreadArgs* ta = t; //convert it into a non unsafe ptr
    (*ta).output = Callgointerpreter((*ta).gof);
}

struct Thread newThread(void* cb) {

    struct ThreadArgs ta;
    ta.gof = cb;
    ta.output = malloc(1); //alloc to the heap

    void* unsafe = &ta;
    HANDLE h = CreateThread(NULL, 0, threadf, unsafe, 0, NULL);

    struct Thread t;
    t.ta = ta;
    t.handle = h;
    return t;
}

void freeInAndOut(struct Thread thread) {
    //used to deallocate the thread, when it goes out of scope
    //but not terminate the thread

    //free the storage
    free(thread.ta.gof);
    free(thread.ta.output);
    //////////////////
}

bool closeThread(struct Thread thread, DWORD exitcode) {

    freeInAndOut(thread);

    return TerminateThread(thread.handle, exitcode);
}

void* waitfor(struct Thread t) {
    WaitForSingleObject(t.handle, INFINITE);
    return t.ta.output;
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