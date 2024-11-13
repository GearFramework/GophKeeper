package model

import (
	"fmt"
	"strings"
	"time"
)

var (
	viewTemplate = `
GUID:        %s
Name:        %s
Description: %s
Date upload: %s
` + strings.Repeat("-", 80)
)

// Entity struct of user entity
type Entity struct {
	GUID        string     `json:"guid" db:"guid"`
	UUID        string     `json:"-"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Type        EntityType `json:"type" db:"type"`
	UploadedAt  time.Time  `json:"uploaded_at" db:"uploaded_at"`
	Attr        MetaData   `json:"attr" db:"attr"`
}

// Entities collection
type Entities = []Entity

// GetFilename return file path of binary data
func (ent *Entity) GetFilename(mountPath string) string {
	if ent.Type == BinaryData {
		return mountPath + "/" + ent.GUID + ent.Attr.Binary.Extension
	}
	return ""
}

// View entity
func (ent *Entity) View() {
	fmt.Printf(viewTemplate,
		ent.GUID,
		ent.Name,
		ent.Description,
		ent.UploadedAt.Format("02/01/2006 15:04"),
	)
	ent.Type.View(ent.Attr)
}
