#ifndef SYSTABLES_SYSCALLS_SELECT_H_
#define SYSTABLES_SYSCALLS_SELECT_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <winsock2.h>
#else
#include <sys/select.h>
#endif

//function to convert raw data (from tusk) to c fdset structures
#define fdset_convert(fdname)                           \
        fd_set fdname;                                  \
        fdname.fd_count = fdname##_count;               \
        /*loop through the fdcount to put sockets*/     \
        for (int i = 0; i < fdname##_count; ++i) {      \
            /*set the current fd socket to the ptr*/    \
            fdname.fd_array[i] =                        \
            (long long int) fdname##_sockets[i];        \
        }                                               \

long long int sysselect(long int nfds, 
    long long int readfds_count, void** readfds_sockets, 
    long long int writefds_count, void** writefds_sockets, 
    long long int exceptfds_count, void** exceptfds_sockets, 
    long long int timeout
) {
    //set the timeouts
    struct timeval tv;
    tv.tv_usec = timeout;

    //convert the fds to fd_sets
    fdset_convert(readfds);
    fdset_convert(writefds);
    fdset_convert(exceptfds);

    return select(nfds, &readfds, &writefds, &exceptfds, &tv);
}

#ifdef __cplusplus
}
#endif

#endif