package auth

import (
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	m := map[string]string{}
	for i := 0; i < 1000000; i++ {
		uuid := CreateUUID("name")
		_, ok := m[uuid]
		assert.False(t, ok)
		m[uuid] = uuid
	}
}

func TestJWT(t *testing.T) {
	if err := logger.Init("info"); err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		uuid := CreateUUID("name")
		tk, err := BuildJWT(uuid)
		assert.NoError(t, err)
		if err != nil {
			continue
		}
		jwtuuid, err := GetUserUUIDFromJWT(tk)
		assert.NoError(t, err)
		if err != nil {
			continue
		}
		assert.Equal(t, uuid, jwtuuid)
		tk += "sds"
		jwtuuid, err = GetUserUUIDFromJWT(tk)
		assert.Error(t, err)
		assert.Equal(t, "", jwtuuid)
	}
}
