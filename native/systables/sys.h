#ifndef TUSK_NATIVE_SYSTABLES_SYS_H_
#define TUSK_NATIVE_SYSTABLES_SYS_H_

#ifdef __cplusplus
extern "C" {
#endif

long int sysread(long int fd, char* buf, unsigned long long int size);
long int syswrite(long int fd, char* buf, unsigned long long int size);
long int sysopen(char* name, int mode);
long int sysclose(int fd);
long int sysstat(char* file);
long int sysfstat(long int fd);
long int syslstat(char* file);

#ifdef __cplusplus
}
#endif

#endif