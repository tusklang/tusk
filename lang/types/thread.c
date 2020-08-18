#ifdef __cplusplus
extern "C" {
#endif

#include "threadwin.h"
#include "threadposix.h"

void freeThread(struct Thread thread) {
    //used to deallocate the thread, when it goes out of scope
    //but not terminate the thread

    //free the storage
    free((*thread.ta).output);
    free(thread.ta);
    //////////////////
}

#ifdef __cplusplus
}
#endif