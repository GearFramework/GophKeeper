package model

import (
	"testing"

	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	tests := []struct {
		et    Entity
		mp    string
		expmp string
	}{
		{et: Entity{Name: "test1", Type: BinaryData, Attr: MetaData{Binary: EntityTypeBinary{Extension: ".jpg"}}}, mp: "/tmp", expmp: "/tmp"},
		{et: Entity{Name: "test2", Type: BinaryData, Attr: MetaData{Binary: EntityTypeBinary{Extension: ".jpg"}}}, mp: "/var", expmp: "/var"},
	}
	for _, test := range tests {
		test.et.GUID = auth.CreateUUID(test.et.Name)
		assert.Equal(t, test.expmp+"/"+test.et.GUID+test.et.Attr.Binary.Extension, test.et.GetFilename(test.mp))
	}
	et := Entity{Name: "test3", Type: PlainText}
	assert.Equal(t, "", et.GetFilename("/tmp"))
}
