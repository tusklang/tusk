package native

import (
	"fmt"
	"os"
)

//OmmPanic panics in an Omm instance
func OmmPanic(err string, line uint64, file string, stacktrace []string) {
	fmt.Println("Panic on line", line, "file", file)
	fmt.Println(err)
	fmt.Println("\nWhen the error was thrown, this was the stack:")
	fmt.Println("  at line", line, "in file", file)
	for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace
		fmt.Println("  " + stacktrace[i])
	}
	os.Exit(1)
}
