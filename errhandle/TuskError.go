package errhandle

//#include "colorprint.h"
import "C"
import (
	"fmt"

	"github.com/tusklang/tusk/tokenizer"
)

type TuskError struct {
	err,
	msg,
	file,
	snippet string
	row, col int
}

func NewTuskError(err, msg, file, snippet string, row, col int) *TuskError {
	return &TuskError{
		err:     err,
		msg:     msg,
		file:    file,
		snippet: snippet,
		row:     row,
		col:     col,
	}
}

func NewTuskErrorFTok(err, msg string, tok tokenizer.Token) *TuskError {
	return NewTuskError(err, msg, tok.File, tok.Snippet, tok.Row, tok.Col)
}

func (e *TuskError) Print() {
	C.errprint(C.CString(
		fmt.Sprintf("error: %s", e.err),
	))
	fmt.Printf("--> at %s:%d:%d\n", e.file, e.row, e.col)
	fmt.Println()
	fmt.Printf("\t%s\n", e.snippet)

	fmt.Print("\t")

	for i := 1; i < e.col-1; i++ {
		fmt.Print(" ")
	}
	fmt.Print("^")

	if e.msg != "" {
		fmt.Printf("-- %s", e.msg)
	}

	fmt.Print("\n\n")
}
