package operations

import "github.com/tusklang/tusk/grouper"

type Operation struct {
	Left  *Operation    //left side operand
	Right *Operation    //right side operand
	Group grouper.Group //operator group/token
}
