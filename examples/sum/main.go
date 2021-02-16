// This is a simple example of a command line tool to sum to values
// go run sum/main.go --a=2 --b=1
// > 2 + 2 = 4

package main

import (
	"fmt"
	"os"

	config "github.com/carlosvin/meta-viper"
)

type cfgStruct struct {
	A float64 `cfg_name:"a" cfg_desc:"Operand A"`
	B float64 `cfg_name:"b" cfg_desc:"Operand B"`
}

func main() {
	cfg := &cfgStruct{
		A: 0,
		B: 0}
	_, err := config.New(cfg, os.Args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%f + %f = %f\n", cfg.A, cfg.B, cfg.A+cfg.B)
}
