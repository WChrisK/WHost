package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"../../internal/appmanager"
	"../../internal/appmanager/zandronum"

	"github.com/julienschmidt/httprouter"
)

var server appmanager.AppServerInfo

func readArgs() *appmanager.RuntimeArgs {
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

	args := &appmanager.RuntimeArgs{}
	if err := json.Unmarshal(decodedData, args); err != nil {
		log.Fatal("Unable to read command line arguments:", err)
	}

	return args
}

func createServer(runtimeArgs *appmanager.RuntimeArgs) {
	createdServer, err := zandronum.CreateServer(runtimeArgs)
	if err != nil {
		log.Fatal("Unable to create server:", err)
	}

	server = createdServer
}

func status(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(server); err != nil {
		log.Println("Error writing status:", err)
	}
}

func runEndpoints() {
	router := httprouter.New()
	router.GET("/status", status)

	log.Fatal(http.ListenAndServe(":8088", nil))
}

func main() {
	runtimeArgs := readArgs()
	createServer(runtimeArgs)
	runEndpoints()
}
