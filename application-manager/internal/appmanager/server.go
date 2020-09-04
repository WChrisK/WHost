package appmanager

type RuntimeArgs struct {
	Args []string `json:"args"`
	Executable string `json:"executable"`
	Files []string `json:"files"`
	Iwad string `json:"iwad"`
	Optfiles []string `json:"optfiles"`
	Wadpath string `json:"wadpath"`
}

type AppServerInfo interface {
	GetCvars() map[string]string
	GetPlayerCount() int
	GetPort() int
}
