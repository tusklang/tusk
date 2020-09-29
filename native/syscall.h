#ifndef TUSK_NATIVE_SYSCALL_H_
#define TUSK_NATIVE_SYSCALL_H_

#ifdef __cplusplus
extern "C" {
#endif

const int MAX_SYS_ARGC = 21;
#define sysproto \
    void* a0,  \
    void* a1,  \
    void* a2,  \
    void* a3,  \
    void* a4,  \
    void* a5,  \
    void* a6,  \
    void* a7,  \
    void* a8,  \
    void* a9,  \
    void* a10, \
    void* a11, \
    void* a12, \
    void* a13, \
    void* a14, \
    void* a15, \
    void* a16, \
    void* a17, \
    void* a18, \
    void* a19

typedef long int (*SYSF)(sysproto);

static inline int makeintfromunsafe(void* v) {
    //prevent the warning, because it works
    #pragma GCC diagnostic ignored "-Wpointer-to-int-cast"
    #pragma GCC diagnostic push
    return (int) v;
    #pragma GCC diagnostic pop
}

static inline void* makeunsafeint(int v) {
    //prevent the warning because it works
    #pragma GCC diagnostic ignored "-Wint-to-pointer-cast"
    #pragma GCC diagnostic push
    return (void*) v;
    #pragma GCC diagnostic pop
}

static inline long int callsys(void* fn, sysproto) {
    long int called = ((SYSF)(fn))(
        a0,
        a1,
        a2,
        a3,
        a4,
        a5,
        a6,
        a7,
        a8,
        a9,
        a10,
        a11,
        a12,
        a13,
        a14,
        a15,
        a16,
        a17,
        a18,
        a19
    ); //call the sycall func
    return called; //return the val
}

#ifdef __cplusplus
}
#endif

#endif