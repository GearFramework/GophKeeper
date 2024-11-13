package client

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		cmd              string
		validExpected    bool
		needAuthExpected bool
	}{
		{"add", true, true},
		{"del", true, true},
		{"list", true, true},
		{"signup", true, false},
		{"view", true, true},
		{"ping", false, true},
	}
	for _, test := range tests {
		cmd := Command(test.cmd)
		assert.Equal(t, test.validExpected, slices.Contains(commmands, cmd))
		assert.Equal(t, test.needAuthExpected, cmd.NeedAuth())
	}
}
