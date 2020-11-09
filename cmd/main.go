package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/carlosvin/go-config-example/internal"
)

type CfgStruct struct {
	Host string `cfg_name:"host" cfg_desc:"Server host"`
	Port int    `cfg_name:"port" cfg_desc:"Server port"`
}

func main() {
	cfg := &CfgStruct{Host: "localhost", Port: 6000}
	config.Load(cfg)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Serving at %v...", cfg)
	http.ListenAndServe(addr, nil)
}
