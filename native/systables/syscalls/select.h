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
#define fdset_convert(fdname)                                               \
    fd_set fdname;                                                          \
    fd_set *fdname##p;                                                      \
    FD_ZERO(&fdname);                                                       \
    if (fdname##_count != 0)                                                \
    {                                                                       \
        /*puts it in reverse but who cares*/                                \
        /*also im cool because i understand post/pre increments and stuff*/ \
        while (--fdname##_count >= 0)                                       \
            FD_SET(                                                         \
                (long long int)                                             \
                    fdname##_sockets[fdname##_count],                       \
                &fdname);                                                   \
        fdname##p = &fdname;                                                \
    }                                                                       \
    else                                                                    \
        fdname##p = NULL;

    long long int sysselect(long int nfds,
                            long long int readfds_count, void **readfds_sockets,
                            long long int writefds_count, void **writefds_sockets,
                            long long int exceptfds_count, void **exceptfds_sockets,
                            long long int timeoutsec, long long int timeoutusec)
    {
        //convert the fds to fd_sets
        fdset_convert(readfds);
        fdset_convert(writefds);
        fdset_convert(exceptfds);

        if (timeoutsec == -1 || timeoutusec == -1)
        {
            int r = select(nfds, readfdsp, writefdsp, exceptfdsp, NULL); //supply a timeout of -1 to have no timeout
            perror("select()");
            return r;
        }

        //set the timeouts
        struct timeval tv;
        tv.tv_sec = timeoutsec;
        tv.tv_usec = timeoutusec;

        int r = select(nfds, readfdsp, writefdsp, exceptfdsp, &tv);
        return r;
    }

#ifdef __cplusplus
}
#endif

#endif