//Package tusksys implements cross platform numeric system calls in tusk based on linux
package tusksys

//#include "sys.h"
import "C"
import "unsafe"

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
}
