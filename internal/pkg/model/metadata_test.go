package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	md := MetaData{Text: EntityTypeText{"a"}}
	dv, err := md.Value()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	mdt := MetaData{}
	err = mdt.Scan(dv)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, md.Text.Value, mdt.Text.Value)
}
