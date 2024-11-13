package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityTypes(t *testing.T) {
	tests := []struct {
		t         string
		expType   EntityType
		expString string
		expValid  bool
	}{
		{t: "credentials", expType: Credentials, expString: "credentials", expValid: true},
		{t: "text", expType: PlainText, expString: "text", expValid: true},
		{t: "binary", expType: BinaryData, expString: "binary", expValid: true},
		{t: "creditcard", expType: Creditcard, expString: "creditcard", expValid: true},
		{t: "media", expType: EntityType("media"), expString: "media", expValid: false},
	}
	for _, test := range tests {
		et := EntityType(test.t)
		assert.Equal(t, test.expType, et)
		assert.Equal(t, test.expString, et.String())
		assert.Equal(t, test.expType.String(), et.String())
		assert.Equal(t, test.expValid, et.Is())
	}
}
