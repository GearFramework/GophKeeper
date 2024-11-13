package server

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		params       []string
		expectedAddr string
	}{
		{params: []string{"-a", ":8090", "-d", "postgres://gk:7766gm@localhost:5432/gk"}, expectedAddr: ":8090"},
		{params: []string{}, expectedAddr: defaultAddress},
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
		fl := ParseFlags()
		conf, err := NewConfig(fl)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedAddr, conf.Addr)
		copy(os.Args, old)
	}

}
