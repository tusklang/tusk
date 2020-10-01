#ifndef SYSTABLES_SYSCALLS_SOCKET_H_
#define SYSTABLES_SYSCALLS_SOCKET_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <winsock.h>
#include <winsock2.h>
#else
#include <sys/socket.h>
#endif

long long int syssocket(int domain, int type, int protocol) {
    return socket(domain, type, protocol);
}

long long int sysconnect(long long int fd, int sa_family, char* sa_data) {
    struct sockaddr addr;
    addr.sa_family = sa_family;
    strcpy(sa_data, sa_data); //use strcpy because arrays are not modifiable
    return connect(fd, &addr, sizeof(addr));
}

long long int sysaccept(long long int fd, int sa_family, char* sa_data) {
    struct sockaddr addr;
    addr.sa_family = sa_family;
    strcpy(sa_data, sa_data); //use strcpy because arrays are not modifiable
    int size = sizeof(addr);
    return accept(fd, &addr, &size);
}

//works for both send and recv
#define sys_recv_impl {                                     \
    struct sockaddr addr;                                   \
    addr.sa_family = sa_family;                             \
    /* use strcpy because arrays are not modifiable */      \
    strcpy(sa_data, sa_data);                               \
    int size = sizeof(addr);                                \
    return sendto(fd, buf, buflen, 0, &addr, sizeof(addr)); \
}

long long int syssendto(long long int fd, char* buf, int buflen, int sa_family, char* sa_data) 
    sys_recv_impl

long long int sysrecvfrom(long long int fd, char* buf, int buflen, int sa_family, char* sa_data) 
    sys_recv_impl

#ifdef __cplusplus
}
#endif

#endif