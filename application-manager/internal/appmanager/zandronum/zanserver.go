package zandronum

import (
	"log"
	"os"
	"os/exec"

	appmanager ".."
)

type ZanServerInfo struct {
	Cvars map[string]string `json:"cvars"`
	PlayerCount int `json:"playerCount"`
	Port int `json:"port"`
}

func (s *ZanServerInfo) GetCvars() map[string]string {
	return s.Cvars
}

func (s *ZanServerInfo) GetPlayerCount() int {
	return s.PlayerCount
}

func (s *ZanServerInfo) GetPort() int {
	return s.Port
}

func CreateServer(args *appmanager.RuntimeArgs) (*ZanServerInfo, error) {
	path, err := exec.LookPath(args.Executable)
	if err != nil {
		log.Fatal("Unable to find executable at: ", path, " - Error: ", err)
	}

	cmd := &exec.Cmd{
		Path: path,
		Args: []string{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin: os.Stdin,
	}

	if err := cmd.Run(); err != nil {
		log.Fatal("Could not run command")
	}

	return &ZanServerInfo{
		Cvars: make(map[string]string),
	}, nil
}
