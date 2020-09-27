#ifndef TUSK_NATIVE_SYSARRAY_H_
#define TUSK_NATIVE_SYSARRAY_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

static inline void** sysarray_make(int size) {
    return (void**) calloc(size, sizeof(void*));
}

static inline void sysarray_setint(void** arr, int index, int val) {
    #pragma GCC diagnostic push 
	#pragma GCC diagnostic ignored "-Wint-to-pointer-cast"
    arr[index] = (void*)val;
    #pragma GCC diagnostic pop 
}

static inline void sysarray_setstr(void** arr, int index, char* val) {
    arr[index] = (void*)val;
}

static inline void* sysarray_get(void** arr, int index) {
    return arr[index];
}

#ifdef __cplusplus
}
#endif

#endif