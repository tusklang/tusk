package types

//package to store all of the datatypes in tusk

type CliParams struct {
	Prec        uint64
	Output      string
	Name        string
	Directory   string
	TuskDirname string
}

//TuskType is an interface representing any arbitrary object in Tusk
type TuskType interface {
	Format() string

	//difference between type and typeof is for prototypes
	//Type() for a proto will return "prototype"
	//TypeOf() will return the proto's name
	Type() string
	TypeOf() string
	//////////////////////////////////////////////////////

	Deallocate()
	Clone() *TuskType
	Range(func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError)
}
