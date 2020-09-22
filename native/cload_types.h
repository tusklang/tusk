#ifndef TUSK_NATIVE_CLOAD_TYPES_H_
#define TUSK_NATIVE_CLOAD_TYPES_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

static inline void** allocarr(unsigned long long size) {
    return (void**) malloc(size);
}

static inline void setindex(void** arr, int index, void* value) {
    arr[index] = value;
}

#ifdef __cplusplus
}
#endif

#endif