package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

// ViewEntity load entity from remote server and show metadata
func (c *GkClient) ViewEntity() error {
	var err error
	guid := ""
	fmt.Print("Enter entity GUID: ")
	_, err = fmt.Fscan(os.Stdin, &guid)
	if err != nil {
		return err
	}
	ent, err := c.getEntity(guid)
	if err != nil {
		return err
	}
	ent.View()
	return nil
}

func (c *GkClient) getEntity(guid string) (*model.Entity, error) {
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"GET", c.Conf.Addr+"/v1/entities/"+guid, nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	var entity model.Entity
	err = json.Unmarshal(body, &entity)
	return &entity, err
}
