package configer

import (
	"flag"
	"io"
	"testing"
)

// TestStringSlice has taken inspiration from Go's `src/flag/flag_test.go` file
// and the function `TestUserDefined`. https://go.dev/src/flag/flag_test.go
func TestStringSlice_multiple(t *testing.T) {
	tts := map[string]struct {
		args         []string
		expectedStr  string
		expectedSize int
	}{
		"empty set": {
			args:         []string{},
			expectedStr:  "",
			expectedSize: 0,
		},
		"single set": {
			args:         []string{"-c", "crash"},
			expectedStr:  "crash",
			expectedSize: 1,
		},
		"double set": {
			args:         []string{"-c", "pulp", "-c", "fiction"},
			expectedStr:  "pulp,fiction",
			expectedSize: 2,
		},
		"triple set": {
			args:         []string{"-c", "v", "-c", "for", "-c", "vendetta"},
			expectedStr:  "v,for,vendetta",
			expectedSize: 3,
		},
		"flag variance set": {
			args:         []string{"-c", "back", "--c", "to", "-c=the", "--c=future"},
			expectedStr:  "back,to,the,future",
			expectedSize: 4,
		},
	}

	for name, tt := range tts {
		t.Run(name, func(t *testing.T) {
			var flags flag.FlagSet
			flags.Init(name, flag.ContinueOnError)
			flags.SetOutput(io.Discard)

			var ss StringSlice
			flags.Var(&ss, "c", "usage")

			if err := flags.Parse(tt.args); err != nil {
				t.Error(err)
			}

			if len(ss) != tt.expectedSize {
				t.Fatalf("len(ss) = %d; expected %d", len(ss), tt.expectedSize)
			}

			if ss.String() != tt.expectedStr {
				t.Errorf("ss.String() = %q; expected: %q", ss.String(), tt.expectedStr)
			}
		})
	}
}
