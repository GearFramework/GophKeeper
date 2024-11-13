package server

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONFlags(t *testing.T) {
	tests := []struct {
		params      []string
		expectedEnv string
	}{
		{params: []string{"-j", "config.json"}, expectedEnv: "config.json"},
		{params: []string{}, expectedEnv: "\x00"},
	}
	var old []string
	copy(old, os.Args)
	for _, test := range tests {
		var buf strings.Builder
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		flag.CommandLine.SetOutput(&buf)
		os.Args = []string{os.Args[0]}
		if len(test.params) > 0 {
			os.Args = append(os.Args, test.params...)
		}
		conf := ParseFlags()
		assert.Equal(t, test.expectedEnv, conf.JSONFile)
		copy(os.Args, old)
	}

}
