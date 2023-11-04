package configer

import (
	"os"
	"strings"
	"testing"
)

type myconfig struct {
	Port    int    `env:"PORT"`
	LogFile string `env:"LOG_FILE"`
}

func (mc myconfig) Validate() error {
	return nil
}

func buildFileWithContents(t *testing.T, pattern, contents string) (string, error) {
	t.Helper()

	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(contents); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func buildConfigFiles(t *testing.T, extensions ...string) ([]string, func(), error) {
	t.Helper()

	files := make([]string, len(extensions))

	for i, ext := range extensions {
		fn, err := func(ext string) (string, error) {
			var pattern, contents string

			switch strings.ToLower(ext) {
			case "json":
				pattern = "config-*.json"
				contents = `{"Port": 3000, "LogFile": "out-json.log"}`
			case "yaml":
				pattern = "config-*.yaml"
				contents = "port: 5000\nlogfile: 'out-yaml.log'"
			case "env":
				pattern = "config-*.env"
				contents = "PORT=8000\nLOG_FILE=out-env.log"
			}

			f, err := os.CreateTemp("", pattern)
			if err != nil {
				return "", err
			}
			defer f.Close()

			// write to file with default data
			_, err = f.WriteString(contents)
			if err != nil {
				return "", err
			}

			return f.Name(), nil
		}(ext)

		if err != nil {
			return nil, nil, err
		}

		files[i] = fn
	}

	cleanup := func() {
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				t.Errorf("unable to remove file %q: %v", f, err)
			}
		}
	}

	return files, cleanup, nil
}

// TestLoad_filename_dot_env tests that the file ".env" is accepted and read
func TestLoad_filename_dot_env(t *testing.T) {
	f, err := os.Create(".env")
	if err != nil {
		t.Fatalf("unable to create .env: %v", err)
	}
	f.WriteString("PORT=8000")
	f.Close()

	t.Cleanup(func() {
		if err := os.Remove(".env"); err != nil {
			t.Errorf("unable to delete file %q: %v", ".env", err)
		}
	})

	mc := myconfig{}
	if err := Load(&mc, ".env"); err != nil {
		t.Errorf("expected nil err; got: %v", err)
	}

	if mc.Port != 8000 {
		t.Errorf("mc.Port not loaded correctly; expected: 8000; got: %v", mc.Port)
	}
}

func TestLoad_invalid_interface(t *testing.T) {
	expected := "invalid struct"
	mc := myconfig{}
	err := Load(mc)
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Errorf("expected err to complain about non-pointer; expected: err contains %q; got %v", expected, err)
	}
}

func TestLoad_json_file(t *testing.T) {
	files, cleanup, err := buildConfigFiles(t, "JSON")
	if err != nil {
		t.Fatalf("unable to create prereq files: %v", err)
	}

	t.Cleanup(cleanup)

	mc := myconfig{Port: 1234}
	if err := Load(&mc, files...); err != nil {
		t.Fatalf("unable to load json file: %v", err)
	}

	// json defaults
	expected := myconfig{Port: 3000, LogFile: "out-json.log"}
	if mc.Port != expected.Port {
		t.Errorf("Load(&mc, tmp.json) load failure; expected mc.Port=%d; got mc.Port=%d", expected.Port, mc.Port)
	}
	if mc.LogFile != expected.LogFile {
		t.Errorf("Load(&mc, tmp.json) load failure; expected mc.LogFile=%s; got mc.LogFile=%s", expected.LogFile, mc.LogFile)
	}
}

func TestLoad_yaml_file(t *testing.T) {
	files, cleanup, err := buildConfigFiles(t, "YAML")
	if err != nil {
		t.Fatalf("unable to create prereq files: %v", err)
	}

	t.Cleanup(cleanup)

	mc := myconfig{Port: 1234}
	if err := Load(&mc, files...); err != nil {
		t.Fatalf("unable to load yaml file: %v", err)
	}

	// yaml defaults
	expected := myconfig{Port: 5000, LogFile: "out-yaml.log"}
	if mc.Port != expected.Port {
		t.Errorf("Load(&mc, tmp.yaml) load failure; expected mc.Port=%d; got mc.Port=%d", expected.Port, mc.Port)
	}
	if mc.LogFile != expected.LogFile {
		t.Errorf("Load(&mc, tmp.yaml) load failure; expected mc.LogFile=%s; got mc.LogFile=%s", expected.LogFile, mc.LogFile)
	}
}

func TestLoad_env_file(t *testing.T) {
	files, cleanup, err := buildConfigFiles(t, "ENV")
	if err != nil {
		t.Fatalf("unable to create prereq files: %v", err)
	}

	t.Cleanup(cleanup)

	mc := myconfig{Port: 1234}
	if err := Load(&mc, files...); err != nil {
		t.Fatalf("unable to load env file: %v", err)
	}

	// env defaults
	expected := myconfig{Port: 8000, LogFile: "out-env.log"}
	if mc.Port != expected.Port {
		t.Errorf("Load(&mc, tmp.env) load failure; expected mc.Port=%d; got mc.Port=%d", expected.Port, mc.Port)
	}
	if mc.LogFile != expected.LogFile {
		t.Errorf("Load(&mc, tmp.env) load failure; expected mc.LogFile=%s; got mc.LogFile=%s", expected.LogFile, mc.LogFile)
	}
}

func TestLoad_multi(t *testing.T) {
	yamlFile, err := buildFileWithContents(t, "config-*.yaml", "logfile: 'out-yaml.log'")
	if err != nil {
		t.Fatalf("unable to create yaml file: %v", err)
	}
	jsonFile, err := buildFileWithContents(t, "config-*.json", `{"Port": 3000}`)
	if err != nil {
		t.Fatalf("unable to create json file: %v", err)
	}

	t.Cleanup(func() {
		for _, f := range []string{yamlFile, jsonFile} {
			if err := os.Remove(f); err != nil {
				t.Errorf("unable to delete file %q: %v", f, err)
			}
		}
	})

	mc := myconfig{}
	if err := Load(&mc, jsonFile, yamlFile); err != nil {
		t.Fatalf("unable to load env file: %v", err)
	}

	expected := myconfig{Port: 3000, LogFile: "out-yaml.log"}
	if mc.Port != expected.Port {
		t.Errorf("Load(&mc, yamlFile, jsonFile) load failure; expected mc.Port=%d; got mc.Port=%d", expected.Port, mc.Port)
	}
	if mc.LogFile != expected.LogFile {
		t.Errorf("Load(&mc, yamlFile, jsonFile) load failure; expected mc.LogFile=%s; got mc.LogFile=%s", expected.LogFile, mc.LogFile)
	}
}

func TestLoad_env_override(t *testing.T) {
	if err := os.Setenv("PORT", "443"); err != nil {
		t.Fatalf("cannot set PORT env: %v", err)
	}

	files, cleanup, err := buildConfigFiles(t, "JSON")
	if err != nil {
		t.Fatalf("unable to create prereq files: %v", err)
	}

	t.Cleanup(cleanup)

	mc := myconfig{Port: 1234}
	if err := Load(&mc, files...); err != nil {
		t.Fatalf("unable to load json file: %v", err)
	}

	if mc.Port != 443 {
		t.Errorf("Load(&mc, files...) env override load failure; expected mc.Port=%d; got mc.Port=%d", 443, mc.Port)
	}
}
