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
long long int sysrealloc(long long int loc, long long int size);
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
long long int syswaitpid(long long int fd, long long int maxtime);
long long int syskillpid(long long int fd, int sig);
long long int sysuname(char* sysname, char* nodename, char* release);
long long int sysfsync(long int fd);
long long int sysftrucate(long int fd, long long int length);
long long int syslsdir(long long int loc, void** subnames, void** fsnumbers, int len);
long long int syssizedir(long long int loc);
long long int sysclosedir(long long int loc);
long long int sysgetcwd(char* buf);
long long int syschdir(char* path);
long long int sysrename(char* oldp, char* newp);
long long int sysmkdir(char* path);
long long int sysrmdir(char* path);
long long int syslink(char* p1, char* p2);
long long int sysunlink(char* path);
long long int syschmod(char* name, int mode);
long long int sysgettime();
long long int sysgettimezone();
long long int syssettime(long long int unixtime);
long long int syssettimezone(long long int lgmt);
long long int syschroot(char* path);
long long int syssync();
long long int sysgethostname(char* name, long long int len);
long long int syssethostname(char* name, long long int len);
long long int sysgetdomainname(char* name, long long int len);
long long int syssetdomainname(char* name, long long int len);
long long int sysgettid();
long long int systkill(long long int tid, int exitc);

#ifdef __cplusplus
}
#endif

#endif