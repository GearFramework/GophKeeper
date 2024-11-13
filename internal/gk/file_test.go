package gk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlags(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"/file.text", false},
		{"./file.go", true},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, IsExistsFile(test.file))
	}
}
