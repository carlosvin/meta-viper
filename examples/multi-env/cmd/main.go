// This is a simple example of a command loading configuration depending on the environment
// e.g. 1: go run ./cmd/main.go --config=qa # it loads configuration from qa.json file
// e.g. 2: go run ./cmd/main.go --config=dev --host=my.local.host # it loads dev configuration and overrides the 'host' value
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	config "github.com/carlosvin/meta-viper"
)

type cfgStruct struct {
	Host      string `cfg_name:"host" cfg_desc:"Server host"`
	Port      int    `cfg_name:"port" cfg_desc:"Server port"`
	SearchAPI string `cfg_name:"apis.search" cfg_desc:"Search API endpoint"`
}

func main() {
	cfg := &cfgStruct{
		Host:      "localhost",
		Port:      6000,
		SearchAPI: "https://google.es"}
	_, err := config.New(cfg, os.Args)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Serving at %v...", cfg)
	http.ListenAndServe(addr, nil)
}
