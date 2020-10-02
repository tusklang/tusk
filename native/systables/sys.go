//Package tusksys implements cross platform numeric system calls in tusk based on linux
package tusksys

import (
	"unsafe"
)

//#cgo windows LDFLAGS: -lwsock32
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
}
