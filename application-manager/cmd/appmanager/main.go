package main

import (
	"../../internal/appmanager"
	"../../internal/appmanager/zandronum"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

var server appmanager.AppServerInfo

func readArgs() *appmanager.ServerRuntimeArgs {
	var data string
	flag.StringVar(&data, "data", "", "The base64 encoded json data")
	flag.Parse()

	if len(data) == 0 {
		log.Fatal("No command line `data` argument provided")
	}

	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatal("Unexpected error when decoding runtime args:", err)
	}

	args := &appmanager.ServerRuntimeArgs{}
	if err := json.Unmarshal(decodedData, args); err != nil {
		log.Fatal("Unable to read command line arguments:", err)
	}

	return args
}

func createServer(runtimeArgs *appmanager.ServerRuntimeArgs) {
	createdServer, err := zandronum.CreateServer(runtimeArgs)
	if err != nil {
		log.Fatal("Unable to create server:", err)
	}

	server = createdServer
}

func status(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(server); err != nil {
		log.Fatal("Error writing status:", err)
	}
}

func shutdown(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	log.Println("Shutdown request received, terminating server")

	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprintf(w, "{\"success\":true}")
	if err != nil {
		log.Println("Error writing shutdown message to host:", err)
	}

	// A hacky and terrible way of making sure the write goes through. This
	// should be done better later on.
	go func() {
		<-(time.NewTimer(3 * time.Second)).C
		os.Exit(0)
	}()
}

func runEndpoints() {
	router := httprouter.New()
	router.GET("/status", status)
	router.POST("/shutdown", shutdown)

	// We only want queries done by our local machine.
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func main() {
	runtimeArgs := readArgs()
	createServer(runtimeArgs)
	runEndpoints()
}
