package client

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	tests := []struct {
		params       []string
		expectedAddr string
		expectedCmd  Command
		expectedUser string
		expectedType string
	}{
		{params: []string{"-a", ":8080", "-c", "signup"}, expectedAddr: ":8080", expectedCmd: CommandSignup, expectedUser: "", expectedType: ""},
		{params: []string{"-a", ":8081", "-c", "add"}, expectedAddr: ":8081", expectedCmd: CommandAdd, expectedUser: "", expectedType: ""},
		{params: []string{"-a", ":8082"}, expectedAddr: ":8082", expectedCmd: "", expectedUser: "", expectedType: ""},
		{params: []string{"-a", ":8083", "-u", "denis"}, expectedAddr: ":8083", expectedCmd: "", expectedUser: "denis", expectedType: ""},
		{params: []string{"-c", "del"}, expectedAddr: "", expectedCmd: CommandDel, expectedUser: "", expectedType: ""},
		{params: []string{"-c", "add", "-t", "binary"}, expectedAddr: "", expectedCmd: CommandAdd, expectedUser: "", expectedType: "binary"},
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
		conf, err := ParseFlags()
		fmt.Println("Config ", conf)
		if err != nil {
			assert.Error(t, ErrInvalidCommand, err)
			copy(os.Args, old)
			continue
		}
		assert.Equal(t, test.expectedAddr, conf.Addr)
		assert.Equal(t, test.expectedUser, conf.Username)
		assert.Equal(t, test.expectedCmd, conf.Command)
		assert.Equal(t, test.expectedType, conf.Type)
		copy(os.Args, old)
	}
}
