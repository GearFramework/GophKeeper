package model

import (
	"fmt"
	"math"
)

var (
	viewCredentialsTemplate = `
Username:    %s
Password:    %s
`
	viewTextTemplate = `
%s
`
	viewBinaryTemplate = `
Filename:    %s
Size:        %s
`
	viewCreditcardTemplate = `
Bank:       %s
Number:     %s
Date:       %d/%d
Cardholder: %s
CVV:        %d
`
)

// EntityType type of entities
type EntityType string

const (
	// Credentials login and password data
	Credentials EntityType = "credentials"
	// PlainText text data and text files
	PlainText EntityType = "text"
	// BinaryData binary files
	BinaryData EntityType = "binary"
	// Creditcard data of credit card
	Creditcard EntityType = "creditcard"
)

// EntityTypeText struct of type text data
type EntityTypeText struct {
	Value string `json:"value"`
}

// EntityTypeCredential struct of type credential
type EntityTypeCredential struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EntityTypeBinary struct of type binary data
type EntityTypeBinary struct {
	OriginalFilename string `json:"filename"`
	Extension        string `json:"extension"`
	Size             int64  `json:"size"`
}

// EntityTypeCreditcard struct of type creditcard data
type EntityTypeCreditcard struct {
	BankName   string `json:"bank_name"`
	Number     string `json:"number"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	CardHolder string `json:"cardholder"`
	CVV        int    `json:"cvv"`
}

// Is return true of type is known
func (et EntityType) Is() bool {
	return et == Creditcard || et == Credentials || et == BinaryData || et == PlainText
}

// String return string equal of type
func (et EntityType) String() string {
	return string(et)
}

// View metadata of type
func (et EntityType) View(md MetaData) {
	switch et {
	case Credentials:
		et.viewCredentials(md.Credential)
	case PlainText:
		et.viewText(md.Text)
	case BinaryData:
		et.viewBinary(md.Binary)
	case Creditcard:
		et.viewCreditcard(md.Creditcard)
	}
}

func (et EntityType) viewCredentials(ent EntityTypeCredential) {
	fmt.Printf(viewCredentialsTemplate+"\n",
		ent.Login,
		ent.Password,
	)
}

func (et EntityType) viewText(ent EntityTypeText) {
	fmt.Printf(viewTextTemplate+"\n",
		ent.Value,
	)
}

func (et EntityType) viewBinary(ent EntityTypeBinary) {
	fmt.Printf(viewBinaryTemplate+"\n",
		ent.OriginalFilename,
		prettyByteSize(ent.Size),
	)
}

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func (et EntityType) viewCreditcard(ent EntityTypeCreditcard) {
	fmt.Printf(viewCreditcardTemplate+"\n",
		ent.BankName,
		ent.Number,
		ent.Month,
		ent.Year,
		ent.CardHolder,
		ent.CVV,
	)
}
