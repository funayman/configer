package configer

import (
	"path/filepath"
	"strings"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

type Config interface {
	Validate() error
}

// LoadConfig takes in multiple supported filepath names and reads the configs
// for them. LoadConfig will _always_ pull from the OS environment variables.
// Providing zero paths will only load from the ENV.
// Supported file types:
// - .env
// - .json
// - .yaml
func Load(c Config, paths ...string) error {
	glc := config.New()

	for _, path := range paths {
		switch ext := strings.ToLower(filepath.Ext(path)); ext {
		case ".json":
			glc.AddFeeder(feeder.Json{Path: path})
		case ".yaml", ".yml":
			glc.AddFeeder(feeder.Yaml{Path: path})
		case ".env":
			glc.AddFeeder(feeder.DotEnv{Path: path})
		}
	}

	// always add in default env; the OS environment variables should always have
	// precedence over config files. According to golobby/config documentation:
	// > Lately added feeders overrides early added ones
	// This is why the Env feeder is added after ranging over the paths
	glc.AddFeeder(feeder.Env{})

	if err := glc.AddStruct(c).Feed(); err != nil {
		return err
	}

	return c.Validate()
}
