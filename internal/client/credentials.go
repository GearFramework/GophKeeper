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

func newEntityCredentials() (*gk.UploadEntityRequest, error) {
	name, err := scanNameEntity()
	if err != nil {
		return nil, err
	}
	attr, err := scanCredentialsEntity()
	if err != nil {
		return nil, err
	}
	ent := gk.UploadEntityRequest{
		Name:        name.Name,
		Description: name.Description,
		Type:        model.Credentials,
		Attr: model.MetaData{
			Credential: *attr,
		},
	}
	return &ent, nil
}

func scanCredentialsEntity() (*model.EntityTypeCredential, error) {
	cred := model.EntityTypeCredential{}
	fmt.Print("Enter credential username: ")
	_, err := fmt.Fscan(os.Stdin, &cred.Login)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter credential password: ")
	_, err = fmt.Fscan(os.Stdin, &cred.Password)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (c *GkClient) addCredentials() error {
	req, err := newEntityCredentials()
	if err != nil {
		return err
	}
	guid, err := c.sendRequestCredentials(req)
	if err != nil {
		return err
	}
	fmt.Println("Created entity DONE; GUID:", guid)
	return nil
}

func (c *GkClient) sendRequestCredentials(ent *gk.UploadEntityRequest) (string, error) {
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
	//fmt.Println("New entity created with GUID: ", string(body))
	return string(body), nil
}
