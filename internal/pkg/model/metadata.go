package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MetaData struct of entity attributes
type MetaData struct {
	Text       EntityTypeText       `json:"text,omitempty" db:"text"`
	Binary     EntityTypeBinary     `json:"binary,omitempty" db:"binary"`
	Credential EntityTypeCredential `json:"credential,omitempty" db:"credential"`
	Creditcard EntityTypeCreditcard `json:"creditcard,omitempty" db:"creditcard"`
}

// Value make MetaData struct implement the driver.Valuer interface
func (m MetaData) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan make MetaData struct implement the sql.Scanner interface
func (m *MetaData) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}
