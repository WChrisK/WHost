package zandronum

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	appmanager ".."
)

type ZanServerInfo struct {
	Ready bool `json:"ready"`
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

func handleStdoutLine(line string, serverInfo *ZanServerInfo) {
	if strings.HasPrefix(line, "IP address ") {
		port, err := strconv.Atoi(strings.Split(line, ":")[1])
		if err != nil {
			log.Fatal("Unexpected port number: ", port)
		}
		serverInfo.Port = port
	} else if strings.HasPrefix(line, "*** ") {
		// As soon as we reach our first map, we consider the server ready to
		// be processed by any queries.
		serverInfo.Ready = true
		// TODO: Not true for `changemap`...
		serverInfo.PlayerCount = 0
	}
}

func stdoutListener(serverInfo *ZanServerInfo, reader io.ReadCloser) {
	bufReader := bufio.NewReader(reader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("Unexpected error reading command stdout pipe:" , err)
			}
		}

		line = strings.TrimSuffix(line, "\n")
		handleStdoutLine(line, serverInfo)
	}

	log.Println("Zandronum server terminated, stdout listener stopping")
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
		Stdin: os.Stdin,
	}

	// Have to grab this before the process starts, it's not allowed to be done
	// after it has started.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error getting process stdout pipe: ", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("Could not run command: ", err)
	}

	zanServerInfo := &ZanServerInfo{
		Cvars: make(map[string]string),
		cmd: cmd,
	}

	go stdoutListener(zanServerInfo, stdout)

	return zanServerInfo, nil
}
