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
	13: C.sysioctl,
	14: C.sysreadv,
	15: C.syswritev,
	16: C.syspipe,
	17: C.sysmalloc,
	18: C.sysrealloc,
	19: C.sysfree,
	20: C.sysselect,
	21: C.sysschedyield,
	22: C.sysdup,
	23: C.sysdup2,
	24: C.syspause,
	25: C.sysgetpid,
	26: C.syssocket,
	27: C.sysconnect,
	28: C.sysaccept,
	29: C.syssendto,
	30: C.sysrecvfrom,
	31: C.sysshutdown,
	32: C.syslisten,
	33: C.sysexecv,
	34: C.sysexit,
	35: C.syswaitpid,
	36: C.syskillpid,
	37: C.sysuname,
	38: C.sysfsync,
	39: C.sysftrucate,
	40: C.syslsdir,
	41: C.syssizedir,
	42: C.sysclosedir, //the meaning on life is "sysclosedir", who would've thought it?
	43: C.sysgetcwd,
	44: C.syschdir,
	45: C.sysrename,
	46: C.sysmkdir,
	47: C.sysrmdir,
	48: C.syslink,
	49: C.sysunlink,
	50: C.syschmod,
	51: C.sysgettime,
	52: C.sysgettimezone,
	53: C.syssettime,
	54: C.syssettimezone,
	55: C.syschroot,
	56: C.syssync,
	57: C.sysgethostname,
	58: C.syssethostname,
	59: C.sysgetdomainname,
	60: C.syssetdomainname,
}
