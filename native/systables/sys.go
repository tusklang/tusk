//Package tusksys implements cross platform numeric system calls in tusk based on linux
package tusksys

import (
	"unsafe"
)

//#include <stdlib.h>
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
	13: C.sysmmap,
	14: C.sysmprotect,
	15: C.sysmunmap,
	16: C.sysioctl,
	17: C.sysreadv,
	18: C.syswritev,
	19: C.syspipe,
	20: C.sysmalloc,
	21: C.sysfree,
	22: C.sysselect,
	23: C.sysshedyield,
}
