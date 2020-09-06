package zandronum

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	appmanager ".."
)

type ZanServerInfo struct {
	Cvars map[string]string `json:"cvars"`
	PlayerCount int `json:"playerCount"`
	Port int `json:"port"`
	cmd *exec.Cmd
}

func addFileArg(argFlag string, wadPath string, files []string, commandLineArgs []string) []string {
	for _, fileName := range files {
		filePath := fileName
		if len(wadPath) != 0 {
			filePath = fmt.Sprintf(`"%s/%s"`, wadPath, fileName)
		}
		commandLineArgs = append(commandLineArgs, argFlag, filePath)
	}

	return commandLineArgs
}

func createCommandLineArgs(args *appmanager.ServerRuntimeArgs) []string {
	var commandLineArgs []string
	commandLineArgs = addFileArg("-iwad", args.WadPath, []string{args.Iwad}, commandLineArgs)
	commandLineArgs = addFileArg("-file", args.WadPath, args.Files, commandLineArgs)
	commandLineArgs = addFileArg("-optfile", args.WadPath, args.Optfiles, commandLineArgs)
	commandLineArgs = append(commandLineArgs, args.Args...)

	return commandLineArgs
}

func CreateServer(args *appmanager.ServerRuntimeArgs) (*ZanServerInfo, error) {
	executablePath, err := exec.LookPath(args.Executable)
	if err != nil {
		log.Fatal("Unable to find executable at: ", executablePath, " - Error: ", err)
	}

	cmd := &exec.Cmd{
		Path: executablePath,
		Dir: filepath.Dir(executablePath),
		Args: createCommandLineArgs(args),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin: os.Stdin,
	}

	// TODO: Wrap reader around stdin/out!

	if err := cmd.Run(); err != nil {
		log.Fatal("Could not run command")
	}

	return &ZanServerInfo{
		Cvars: make(map[string]string),
		cmd: cmd,
	}, nil
}
