package appmanager

type ServerRuntimeArgs struct {
	Args       []string `json:"args"`
	Executable string   `json:"executable"`
	Files      []string `json:"files"`
	Iwad       string   `json:"iwad"`
	Optfiles   []string `json:"optfiles"`
	WadPath    string   `json:"wadpath"`
}

type AppServerInfo interface {
}
