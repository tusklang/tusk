#ifndef _WIN32

#ifndef TUSK_NATIVE_SYSCALL_UNIX_H_
#define TUSK_NATIVE_SYSCALL_UNIX_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <errno.h>
#include <stdbool.h>

static inline long int tusksyscall(void* a0, void* a1, void* a2, void* a3) {
	return syscall(*((double*) a0), a1, a2, a3);
}

#ifdef __cplusplus
}
#endif

#endif

#endif