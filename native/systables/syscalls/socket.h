#ifndef SYSTABLES_SYSCALLS_SOCKET_H_
#define SYSTABLES_SYSCALLS_SOCKET_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#include <winsock.h>
#include <winsock2.h>
#else
#include <netdb.h>
#include <sys/socket.h>
#endif

//def to set the socket address struct from raw tusk input
#define setaddr \
    struct sockaddr_in addr_in;                             \
    addr_in.sin_family = sin_family;                        \
    addr_in.sin_addr.s_addr =                               \
    inet_addr( /* set s_addr based on hostname */           \
        inet_ntoa(                                          \
            *(struct in_addr*)(*gethostbyname(hostname)     \
            ->h_addr_list)                                  \
        )                                                   \
    );                                                      \
    addr_in.sin_port = htons(port);                         \
    int size = sizeof(addr_in);                             \
    struct sockaddr* addr = (struct sockaddr*) &addr_in;    \

long long int syssocket(int domain, int type, char* protocol, int port) {
    struct protoent* protoent = getprotobyname(protocol);
    return socket(domain, type, protoent->p_proto);
}

long long int sysconnect(long long int fd, int sin_family, char* hostname, int port) {
    setaddr;
    return connect(fd, addr, size);
}

long long int sysaccept(long long int fd, int sin_family, char* hostname, int port) {
    setaddr;
    return accept(fd, addr, &size);
}

//works for both send and recv
#define sys_recv_impl {                                     \
    setaddr;                                                \
    return sendto(fd, buf, buflen, 0, addr, sizeof(addr)); \
}

long long int syssendto(long long int fd, char* buf, int buflen, int sin_family, char* hostname, int port) 
    sys_recv_impl

long long int sysrecvfrom(long long int fd, char* buf, int buflen, int sin_family, char* hostname, int port) 
    sys_recv_impl

long long int sysshutdown(long int fd, int how) {
    return shutdown(fd, how);
}

long long int sysbind(long int fd, int sin_family, char* hostname, int port) {
    setaddr;
    return bind(fd, addr, size);
}

long long int syslisten(long int fd, int backlog) {
    return listen(fd, backlog);
}

#ifdef __cplusplus
}
#endif

#endif