#ifndef OMM_POSIX_THREAD_H
#define OMM_POSIX_THREAD_H

#ifdef __cplusplus
extern "C" {
#endif

//posix (unix/linux) platforms
#ifndef _WIN32

#include <stdlib.h>
#include <stdbool.h>
#include <pthread.h>
#include "thread.h"

void* threadf(void* tav) {
    struct ThreadArgs* ta = ((struct ThreadArgs*) tav);
    CallGoCB((*ta).gof, (*ta).output);
    pthread_exit(NULL);
}

struct Thread newThread(unsigned long long ptr) {

    struct ThreadArgs ta;
    ta.gof = ptr;
    ta.output = (void**) malloc(1); //alloc to the heap

    struct ThreadArgs* taheap = (struct ThreadArgs*) calloc(1, sizeof(struct ThreadArgs));
    *taheap = ta;

    pthread_t thread;
    pthread_create(&thread, NULL, &threadf, (void*) taheap);

    struct Thread t;
    t.ta = taheap;
    t.handle = thread;
    return t;
}

void* waitfor(struct Thread t) {
    pthread_join(t.handle, NULL);
    return *(*t.ta).output;
}

#endif

#ifdef __cplusplus
}
#endif

#endif
