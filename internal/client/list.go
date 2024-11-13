package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
)

const (
	headerLine = "GUID                                     Name                 Description"
	splitLine  = "---------------------------------------- -------------------- ----------------------------------------"
)

func (c *GkClient) listEntities() error {
	resp, err := c.getList()
	if err != nil {
		return err
	}
	fmt.Printf("\nFound entites %d:\n\n", resp.Count)
	fmt.Println(headerLine)
	fmt.Println(splitLine)
	for _, ent := range *resp.Items {
		spaceNameLen := 0
		name := ent.Name
		if len(ent.Name) < 20 {
			spaceNameLen = 20 - len(name)
		} else {
			name = ent.Name[0:17] + "..."
		}
		spaceDescLen := 0
		desc := ent.Description
		if len(ent.Name) < 40 {
			spaceDescLen = 40 - len(desc)
		} else {
			name = ent.Name[0:37] + "..."
		}
		fmt.Printf("%s"+strings.Repeat(" ", 4)+" %s"+strings.Repeat(" ", spaceNameLen)+" %s"+strings.Repeat(" ", spaceDescLen)+"\n",
			ent.GUID,
			name,
			desc,
		)
	}
	fmt.Println(splitLine)
	return nil
}

func (c *GkClient) getList() (*gk.ListEntitiesResponse, error) {
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"GET", c.Conf.Addr+"/v1/entities", nil,
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
	var list gk.ListEntitiesResponse
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}
