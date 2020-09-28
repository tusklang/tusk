//Package tusksys implements cross platform numeric system calls in tusk based on linux
package tusksys

//#include "sys.h"
import "C"

//SysTable represents all system calls available in tusk
var SysTable = map[int]C.SYSF{
	0: C.SYSF(C.sysread),
	1: C.SYSF(C.syswrite),
	2: C.SYSF(C.sysopen),
	3: C.SYSF(C.sysclose),
}
