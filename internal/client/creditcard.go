package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

func newEntityCreditcard() (*gk.UploadEntityRequest, error) {
	name, err := scanNameEntity()
	if err != nil {
		return nil, err
	}
	attr, err := scanCreditcardEntity()
	if err != nil {
		return nil, err
	}
	ent := gk.UploadEntityRequest{
		Name:        name.Name,
		Description: name.Description,
		Type:        model.Creditcard,
		Attr: model.MetaData{
			Creditcard: *attr,
		},
	}
	return &ent, nil
}

func scanCreditcardEntity() (*model.EntityTypeCreditcard, error) {
	cred := model.EntityTypeCreditcard{}
	fmt.Print("Enter bank name: ")
	_, err := fmt.Fscan(os.Stdin, &cred.BankName)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter number: ")
	_, err = fmt.Fscan(os.Stdin, &cred.Number)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter month: ")
	_, err = fmt.Fscan(os.Stdin, &cred.Month)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter year: ")
	_, err = fmt.Fscan(os.Stdin, &cred.Year)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter cardholder: ")
	_, err = fmt.Fscan(os.Stdin, &cred.CardHolder)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter CVV: ")
	_, err = fmt.Fscan(os.Stdin, &cred.CVV)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (c *GkClient) addCreditcard() error {
	req, err := newEntityCreditcard()
	if err != nil {
		return err
	}
	guid, err := c.sendRequestCreditcard(req)
	if err != nil {
		return err
	}
	fmt.Println("Created entity DONE; GUID:", guid)
	return nil
}

func (c *GkClient) sendRequestCreditcard(ent *gk.UploadEntityRequest) (string, error) {
	b, err := json.Marshal(ent)
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"POST", c.Conf.Addr+"/v1/entities", bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", NewError(resp.StatusCode)
	}
	return string(body), nil
}
