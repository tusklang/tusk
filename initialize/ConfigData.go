package initialize

type ConfigData struct {
	Entry        string   `json:"entry"`
	Dependencies []string `json:"dependencies"`
}
