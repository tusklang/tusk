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
long long int sysmmap(void* addr, unsigned long long int length, int prot, int flags, int fd, long int offset);
long long int sysmprotect(long long int addr, int dwSize, long long int flNewProtect, long int lpflOldProtect);

#ifdef __cplusplus
}
#endif

#endif