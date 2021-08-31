package initialize

type File struct {
	Name      string
	Public    []Declaration
	Protected []Declaration
	Private   []Declaration
}
