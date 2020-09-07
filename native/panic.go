package native

import (
	"errors"
	"fmt"
	"os"
)

//MakeOmmPanic generates the message given in an Omm Panic
func MakeOmmPanic(err string, line uint64, file string, stacktrace []string) error {
	var final string
	final += fmt.Sprintln("Panic on line", line, "file", file)
	final += err
	final += "\nWhen the error was thrown, this was the stack:\n"
	final += fmt.Sprint("  at line ", line, " in file ", file) + "\n"
	for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace

		endl := "\n"
		if i == 0 {
			endl = ""
		}

		final += "  " + stacktrace[i] + endl
	}
	return errors.New(final)
}

//OmmPanic panics in an Omm instance
func OmmPanic(err string, line uint64, file string, stacktrace []string) {
	fmt.Println(MakeOmmPanic(err, line, file, stacktrace))
	os.Exit(1)
}
