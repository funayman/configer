package main

import (
	"flag"
	"fmt"

	"github.com/funayman/configer"
)

type appconfig struct {
	Name string `env:"NAME"`
	Log  struct {
		Level string `env:"LOG_LEVEL"`
		File  string `env:"LOG_FILE"`
	}
	Environment string `env:"ENVIRONMENT"`
	IDs         []int  `env:"IDS"`
}

// Validate satisfies the configer.Config interface
// for examples purposes, returning nil.
func (ac appconfig) Validate() error {
	return nil
}

var (
	configFiles configer.StringSlice
	cfg         appconfig
)

func init() {
	flag.Var(&configFiles, "config", "loads config files")
	flag.Parse()
}

func main() {
	if err := configer.Load(&cfg, configFiles...); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}
