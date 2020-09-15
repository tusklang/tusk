//+build windows

package interpreter

import (
	"fmt"
	"syscall"
)

func makeSyscall(argc int, argv [19]uintptr) uintptr {

	fmt.Println(argc)

	c, _, _ := syscall.Syscall18(
		argv[0],
		uintptr(argc),
		argv[1],
		argv[2],
		argv[3],
		argv[4],
		argv[5],
		argv[6],
		argv[7],
		argv[8],
		argv[9],
		argv[10],
		argv[11],
		argv[12],
		argv[13],
		argv[14],
		argv[15],
		argv[16],
		argv[17],
		argv[18],
	)

	return c
}
