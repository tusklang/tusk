#ifndef TUSK_NATIVE_SYSTABLES_SYS_H_
#define TUSK_NATIVE_SYSTABLES_SYS_H_

#ifdef __cplusplus
extern "C" {
#endif

long int sysread(long int fd, char* buf, unsigned long long int size);
long int syswrite(long int fd, char* buf, unsigned long long int size);
long int sysopen(char* name, int mode);
long int sysclose(int fd);
long int fst_dev(long int fd);
long int fst_ino(long int fd);
long int fst_mode(long int fd);
long int fst_nlink(long int fd);
long int fst_uid(long int fd);
long int fst_gid(long int fd);
long int fst_rdev(long int fd);
long int fst_size(long int fd);

#ifdef __cplusplus
}
#endif

#endif