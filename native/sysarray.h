#ifndef TUSK_NATIVE_SYSARRAY_H_
#define TUSK_NATIVE_SYSARRAY_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

static inline void** sysarray_make(int size) {
    return (void**) malloc(size);
}

static inline void sysarray_set(void** arr, int index, void* val) {
    arr[index] = val;
}

static inline void* sysarray_get(void** arr, int index) {
    return arr[index];
}

#ifdef __cplusplus
}
#endif

#endif