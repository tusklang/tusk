package native

//ErrCodes is a list of the standard runtime error codes in tusk
var ErrCodes = map[string]int{
	"IOB":               0,  //index out of bounds
	"DBZ":               1,  //divide by zero error
	"NOMEM":             2,  //no memory
	"LACKPERM":          3,  //lacking permissions
	"CORRUPTOAT":        4,  //oat file is corrupted
	"SIGNOMATCH":        5,  //function signature does not match
	"ITEMNOTFOUND":      6,  //a hash, proto, or object does not contain an item
	"INVALIDLIT":        7,  //invalid literal
	"UNCLONEABLE":       8,  //type cannot be cloned
	"INVALIDARG":        9,  //invalid argument
	"INVALIDSYSNO":      10, //invalid sysno
	"INVALIDCAST":       11, //invalid typecast
	"NILPTR":            12, //nil pointer reference
	"OPNOTFOUND":        13, //operation signature not found
	"FDWAITERR":         14, //error while selecting on a fd
	"DATANOTRECV":       15, //data not recieved
	"SOCKCANNOTCONNECT": 16, //socket cannot be connected
	"SOCKCANNOTWRITE":   17, //socket cannot be written to
}
