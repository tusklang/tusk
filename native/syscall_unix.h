#ifndef _WIN32 //unix is easy

#ifndef TUSK_NATIVE_SYSCALL_UNIX_H_
#define TUSK_NATIVE_SYSCALL_UNIX_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>
#include <sys/syscall.h>
#include <errno.h>
#include <stdbool.h>

static inline long int tusksyscall(void* a0, void* a1, void* a2, void* a3) {
	/* avoid the warning (because it works) */
	#pragma GCC diagnostic push 
	#pragma GCC diagnostic ignored "-Wpointer-to-int-cast"
	return syscall((int)a0, a1, a2, a3);
	#pragma GCC diagnostic push
}

#ifdef __cplusplus
}
#endif

#endif

#endif