package parser

type ConfigData struct {
	Entry        string   `json:"entry"`
	Dependencies []string `json:"dependencies"`
}
