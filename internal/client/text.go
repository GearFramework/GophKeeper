package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

func newEntityText() (*gk.UploadEntityRequest, error) {
	name, err := scanNameEntity()
	if err != nil {
		return nil, err
	}
	attr, err := scanTextEntity()
	if err != nil {
		return nil, err
	}
	ent := gk.UploadEntityRequest{
		Name:        name.Name,
		Description: name.Description,
		Type:        model.PlainText,
		Attr: model.MetaData{
			Text: *attr,
		},
	}
	return &ent, nil
}

func scanTextEntity() (*model.EntityTypeText, error) {
	txt := model.EntityTypeText{}
	fmt.Print("Enter text: ")
	scan := func() string {
		reader := bufio.NewReader(os.Stdin)
		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			if len(strings.TrimSpace(line)) == 0 {
				break
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n")
	}
	txt.Value = scan()
	return &txt, nil
}

func (c *GkClient) addText() error {
	req, err := newEntityText()
	if err != nil {
		return err
	}
	guid, err := c.sendRequestText(req)
	if err != nil {
		return err
	}
	fmt.Println("Created entity DONE; GUID:", guid)
	return nil
}

func (c *GkClient) sendRequestText(ent *gk.UploadEntityRequest) (string, error) {
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
