#ifndef SYSTABLES_SYSCALLS_SELECT_H_
#define SYSTABLES_SYSCALLS_SELECT_H_

#ifdef __cplusplus
extern "C"
{
#endif

#ifdef _WIN32
#include <winsock2.h>
#else
#include <sys/select.h>
#endif

//function to convert raw data (from tusk) to c fdset structures
#define fdset_convert(fdname)                         \
    fd_set *fdname = NULL;                            \
    if (fdname##_count != 0)                          \
    {                                                 \
        fdname = calloc(1, sizeof(fd_set));           \
        FD_ZERO(fdname);                              \
        for (; fdname##_count >= 0; fdname##_count--) \
            FD_SET(                                   \
                (long long int)                       \
                    fdname##_sockets[fdname##_count], \
                fdname);                              \
    }

    long long int sysselect(long int nfds,
                            long long int readfds_count, void **readfds_sockets,
                            long long int writefds_count, void **writefds_sockets,
                            long long int exceptfds_count, void **exceptfds_sockets,
                            long long int timeoutsec, long long int timeoutusec)
    {
        //set the timeouts
        struct timeval tv;
        tv.tv_sec = timeoutsec;
        tv.tv_usec = timeoutusec;

        //convert the fds to fd_sets
        fdset_convert(readfds);
        fdset_convert(writefds);
        fdset_convert(exceptfds);

        long long int r;

        if (timeoutsec == -1 || timeoutusec == -1)
            r = select(nfds, readfds, writefds, exceptfds, NULL); //supply a timeout of -1 to have no timeout
        r = select(nfds, readfds, writefds, exceptfds, &tv);

        if (readfds != NULL)
            free(readfds);
        if (writefds != NULL)
            free(writefds);
        if (exceptfds != NULL)
            free(exceptfds);

        return r;
    }

#ifdef __cplusplus
}
#endif

#endif