//Package tusksys implements cross platform numeric system calls in tusk based on linux
package tusksys

import (
	"unsafe"
)

//#cgo windows LDFLAGS: -lwsock32 -lkernel32 -lpsapi
//#include "sys.h"
import "C"

//SysTable represents all system calls available in tusk
var SysTable = map[int]unsafe.Pointer{
	0:  C.sysread,
	1:  C.syswrite,
	2:  C.sysopen,
	3:  C.sysclose,
	4:  C.fst_dev,
	5:  C.fst_ino,
	6:  C.fst_mode,
	7:  C.fst_nlink,
	8:  C.fst_uid,
	9:  C.fst_gid,
	10: C.fst_rdev,
	11: C.fst_size,
	12: C.syslseek,
	13: C.sysreadv,
	14: C.syswritev,
	15: C.syspipe,
	16: C.sysmalloc,
	17: C.sysrealloc,
	18: C.sysfree,
	19: C.sysselect,
	20: C.sysschedyield,
	21: C.sysdup,
	22: C.sysdup2,
	23: C.syspause,
	24: C.sysgetpid,
	25: C.syssocket,
	26: C.sysconnect,
	27: C.sysaccept,
	28: C.syssendto,
	29: C.sysrecvfrom,
	30: C.sysshutdown,
	31: C.syslisten,
	32: C.sysexecv,
	33: C.sysexit,
	34: C.syswaitpid,
	35: C.syskillpid,

	//add a syscall 36 later, i removed one of them that was useless, but the lack of a 36 triggers me

	37: C.sysfsync,
	38: C.sysftrucate,
	39: C.syslsdir,
	40: C.syssizedir,
	41: C.sysclosedir, //the meaning on life is "sysclosedir", who would've thought it?
	42: C.sysgetcwd,
	43: C.syschdir,
	44: C.sysrename,
	45: C.sysmkdir,
	46: C.sysrmdir,
	47: C.syslink,
	48: C.sysunlink,
	49: C.syschmod,
	50: C.sysgettime,
	51: C.sysgettimezone,
	52: C.syssettime,
	53: C.syssettimezone,
	54: C.syschroot,
	55: C.syssync,
	56: C.sysgethostname,
	57: C.syssethostname,
	58: C.sysgetdomainname,
	59: C.syssetdomainname,
	60: C.sysgettid,
	61: C.systkill,
}
