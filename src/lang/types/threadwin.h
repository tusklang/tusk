#ifndef OMM_WIN_THREAD_H
#define OMM_WIN_THREAD_H

#ifdef __cplusplus
extern "C" {
#endif

//win32 platforms
#ifdef _WIN32

#include <stdlib.h>
#include <windows.h>
#include <stdbool.h>
#include "thread.h"

DWORD WINAPI threadf(LPVOID tav) {
    struct ThreadArgs* ta = ((struct ThreadArgs*) tav);
    CallGoCB((*ta).gof, (*ta).output);
}

struct Thread newThread(unsigned long long ptr) {

    struct ThreadArgs ta;
    ta.gof = ptr;
    ta.output = (void**) malloc(1); //alloc to the heap

    struct ThreadArgs* taheap = (struct ThreadArgs*) calloc(1, sizeof(struct ThreadArgs));
    *taheap = ta;

    HANDLE h = CreateThread(NULL, 0, threadf, (LPVOID) taheap, 0, NULL);

    struct Thread t;
    t.ta = taheap;
    t.handle = h;
    return t;
}

void* waitfor(struct Thread t) {
    WaitForSingleObject(t.handle, INFINITE);
    return *(*t.ta).output;
}

#endif

#ifdef __cplusplus
}
#endif

#endif
