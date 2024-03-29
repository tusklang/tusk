package errhandle

//#include "colorprint.h"
import "C"
import (
	"fmt"
	"os"
	"strconv"

	"github.com/dlclark/regexp2"
	"github.com/tusklang/tusk/tokenizer"
)

type TuskError struct {
	err,
	msg,
	file,
	snippet string
	row, col int

	/*
		0 - compile
		1 - parse
	*/
	typ int

	strrow string //row as a string
}

func NewCompileError(err, msg, file, snippet string, row, col int) *TuskError {
	return &TuskError{
		err:     err,
		msg:     msg,
		file:    file,
		snippet: snippet,
		row:     row,
		col:     col,
		typ:     0,
	}
}

func NewCompileErrorFTok(err, msg string, tok tokenizer.Token) *TuskError {
	return NewCompileError(err, msg, tok.File, tok.Snippet, tok.Row, tok.Col)
}

func NewParseError(err, msg, file, snippet string, row, col int) *TuskError {
	return &TuskError{
		err:     err,
		msg:     msg,
		file:    file,
		snippet: snippet,
		row:     row,
		col:     col,
		typ:     1,
	}
}

func NewParseErrorFTok(err, msg string, tok tokenizer.Token) *TuskError {
	return NewParseError(err, msg, tok.File, tok.Snippet, tok.Row, tok.Col)
}

func (e *TuskError) printlinepad(printlinen bool /*to print the line number or not*/) {
	//this method prints the pipes at the beginning
	//e.g.
	//	error: whatever
	//	--> at dummy.tusk:0:0
	//	  |
	//	3 |		var hi: i32 = "hello";
	//	  |
	//	  ^ these pipes

	if e.strrow == "" {
		e.strrow = strconv.Itoa(e.row)
	}

	if !printlinen {
		for i := 0; i < len(e.strrow); i++ {
			fmt.Fprint(os.Stderr, " ")
		}
	} else {
		fmt.Fprint(os.Stderr, e.strrow)
	}

	fmt.Fprint(os.Stderr, " |")
}

//print the "errortype: msg" part
func (e *TuskError) printNotice() {
	noti := C.CString(e.err)

	switch e.typ {
	case 0:
		//compile-time error
		C.compileErrorPrint(noti)
	case 1:
		//parse-time error
		C.parseErrorPrint(noti)
	}
}

func (e *TuskError) Print() {
	e.printNotice()
	fmt.Fprintf(os.Stderr, "--> at %s:%d:%d\n", e.file, e.row, e.col)
	e.printlinepad(false)
	fmt.Fprintln(os.Stderr)

	whitespaceR := regexp2.MustCompile("\\s", 0)

	colp := e.col
	snipp := e.snippet

	//remove leading whitespace from the string
	//(yes i know strings.TrimSpace but we also need to count it to subtract it from the col)
	//(because when we print the arrow, we rely on col)
	//(and if col isn't subtracted from then we'd be printing farther away)
	if len(snipp) != 0 {
		for t, _ := whitespaceR.MatchString(string(snipp[0])); t; t, _ = whitespaceR.MatchString(string(snipp[0])) {
			snipp = snipp[1:]
			colp--
		}
	}

	e.printlinepad(true)
	fmt.Fprintf(os.Stderr, "\t%s\n", snipp)

	e.printlinepad(false)
	fmt.Fprint(os.Stderr, "\t")

	for i := 1; i < colp; i++ {
		fmt.Fprint(os.Stderr, " ")
	}
	fmt.Fprint(os.Stderr, "^")

	if e.msg != "" {
		fmt.Fprintf(os.Stderr, "-- %s", e.msg)
	}

	fmt.Fprint(os.Stderr, "\n")
	e.printlinepad(false)
	fmt.Fprint(os.Stderr, "\n")
}
