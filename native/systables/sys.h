#ifndef TUSK_NATIVE_SYSTABLES_SYS_H_
#define TUSK_NATIVE_SYSTABLES_SYS_H_

#ifdef __cplusplus
extern "C" {
#endif

long long int sysread(long int fd, char* buf, unsigned long long int size);
long long int syswrite(long int fd, char* buf, unsigned long long int size);
long long int sysopen(char* name, int mode);
long long int sysclose(int fd);
long long int fst_dev(long int fd);
long long int fst_ino(long int fd);
long long int fst_mode(long int fd);
long long int fst_nlink(long int fd);
long long int fst_uid(long int fd);
long long int fst_gid(long int fd);
long long int fst_rdev(long int fd);
long long int fst_size(long int fd);
long long int syslseek(long int fd, long int offset, int whence);
long long int sysioctl(long long int fd, long long int request, char* argp);
long long int sysreadv(int fd, void** iov_bases, void** iov_lens, int iovcnt);
long long int syswritev(int fd, void** iov_bases, void** iov_lens, int iovcnt);
long long int syspipe(void** fds, long long int size);
long long int sysmalloc(long long int size);
long long int sysfree(long long int ptr);
long long int sysselect(long int nfds, 
    long long int readfds_count, void** readfds_sockets, 
    long long int writefds_count, void** writefds_sockets, 
    long long int exceptfds_count, void** exceptfds_sockets, 
    long long int timeout
);
long long int sysschedyield();
long long int sysdup(long long int fd);
long long int sysdup2(long long int fd, long long int nfd);
long long int syspause();
long long int sysgetpid();
long long int syssocket(int domain, int type, int protocol);
long long int sysconnect(long long int fd, int sa_family, char* sa_data);
long long int sysaccept(long long int fd, int sa_family, char* sa_data);
long long int syssendto(long long int fd, char* buf, int buflen, int sa_family, char* sa_data);
long long int sysrecvfrom(long long int fd, char* buf, int buflen, int sa_family, char* sa_data);
long long int sysshutdown(long int fd, int how);
long long int syslisten(long int fd, int backlog);
long long int sysexecv(char* path, void** argv, void** newenviron);
long long int sysexit(int code);

#ifdef __cplusplus
}
#endif

#endif