#ifndef TUSK_NATIVE_ARRAYC_H_
#define TUSK_NATIVE_ARRAYC_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

static inline void** makecarray(long int len) {
    return (void**) calloc(len, sizeof(void*));
}

static inline void setcarray(void** arr, int idx, void* val) {
    arr[idx] = val;
}

static inline void* getcarray(void** arr, int idx) {
    return arr[idx];
}

#ifdef __cplusplus
}
#endif

#endif