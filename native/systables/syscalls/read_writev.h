#ifndef SYSTABLES_SYSCALLS_READ_WRITEV_H_
#define SYSTABLES_SYSCALLS_READ_WRITEV_H_

#ifdef __cplusplus
extern "C" {
#endif

#define sysreadv_sig long long int sysreadv(int fd, void** iov_bases, void** iov_lens, int iovcnt)
#define syswritev_sig long long int syswritev(int fd, void** iov_bases, void** iov_lens, int iovcnt)

#ifdef _WIN32
//do windows later

#include "read.h"
#include "write.h"

//readv and writev implementation (works for both)
//the empty /**/ are just line breaks
#define read_writev_impl(fn)                                                               \
    for (int i = 0; i < iovcnt; ++i)                                                       \
        if (fn(fd, (char*)(iov_bases[i]), (long long int)(iov_lens[i])) == -1) return -1;  \
    return 0;

sysreadv_sig {
    read_writev_impl(sysread)
}

syswritev_sig {
    read_writev_impl(syswrite)
}

#else

#include <sys/uio.h>
#include <stdlib.h>
#include <string.h>

//readv and writev implementation (works for both)
//the empty /**/ are just line breaks
#define read_writev_impl(fn)                                    \
    /* clone it into the iovec */                               \
    struct iovec* iovec = calloc(iovcnt, sizeof(struct iovec)); \
    int i;                                                      \
    for (i = 0; i < iovcnt; ++i) {                              \
        struct iovec cur;                                       \
        cur.iov_base = iov_bases[i];                            \
        cur.iov_len = (long long int) (iov_lens[i]);            \
        iovec[i] = cur;                                         \
    }                                                           \
    /**/                                                        \
    /* make the syscall */                                      \
    int ret = fn(fd, iovec, iovcnt);                            \
    /**/                                                        \
    /* pass back the iov bases (only for readv) */              \
    if (strcmp(#fn, "readv") == 0)                              \
        for (i = 0; i < iovcnt; ++i)                            \
            iov_bases[i] = iovec[i].iov_base;                   \
    /**/                                                        \
    /* cleanup the iovec */                                     \
    free(iovec);                                                \
    /**/                                                        \
    return ret;                                                 \

sysreadv_sig {
    read_writev_impl(readv);    
}

syswritev_sig {
    read_writev_impl(writev);    
}

#endif

#ifdef __cplusplus
}
#endif

#endif
