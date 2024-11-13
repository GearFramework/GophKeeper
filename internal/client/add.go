package client

import (
	"errors"

	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

var (
	// ErrUnknownDataType unknown data type of new entity
	ErrUnknownDataType = errors.New("unknown data type")
)

// Add new entity into remote server
func (c *GkClient) add() error {
	t := model.EntityType(c.Conf.Type)
	if t == model.Credentials {
		return c.addCredentials()
	}
	if t == model.PlainText {
		return c.addText()
	}
	if t == model.Creditcard {
		return c.addCreditcard()
	}
	if t == model.BinaryData {
		return c.addBinary()
	}
	return ErrUnknownDataType
}
