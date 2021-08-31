package initialize

type Package struct {

	/*
		`Name` is an array, rather than a string, because of sub-directories
		example:

		util/
			arrays/
				ArrayCounter.tusk
		Main.tusk
		tusk.config.json

		Main.tusk would reference the ArrayCoutner methods with `util.arrays`
		The arrays package would have the name []string{"util", "arrays"}
	*/
	Name  []string
	Files []*File
}
