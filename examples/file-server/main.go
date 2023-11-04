package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/funayman/configer"
)

// appconfig is a simple example of what a configuration might look like for a
// server application. The `env` tags must be provided to be picked up by the
// OS's environment variables.
// To learn more refer to [golobby/env](https://github.com/golobby/env)
type appconfig struct {
	Server struct {
		Addr string `env:"SERVER_ADDR"`
		Port int    `env:"SERVER_PORT"`
	}
	Directory string `env:"DIRECTORY"`
}

// Validate satisfies the configer.Config interface
func (ac appconfig) Validate() error {
	switch "" {
	case ac.Server.Addr:
		return fmt.Errorf("missing required config variable %q", "SEVER_ADDR")
	case ac.Directory:
		return fmt.Errorf("missing required config variable %q", "DIRECTORY")
	}

	if ac.Server.Port <= 0 || ac.Server.Port > 65535 {
		return fmt.Errorf("%q must be between 1 and 65535", "SERVER_PORT")
	}

	return nil
}

var (
	cfg appconfig
)

func main() {
	if err := configer.Load(&cfg, "config.json"); err != nil {
		panic(err)
	}

	// serve directory provided by config
	http.Handle("/", http.FileServer(http.Dir(cfg.Directory)))

	addr := net.JoinHostPort(
		cfg.Server.Addr,
		fmt.Sprintf("%d", cfg.Server.Port), // super lazy
	)
	fmt.Printf("serving directory %q on %s\n", cfg.Directory, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
