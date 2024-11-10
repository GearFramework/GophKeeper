package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

func newBinaryEntity() (*gk.UploadEntityRequest, string, error) {
	name, err := scanNameEntity()
	if err != nil {
		return nil, "", err
	}
	attr, fp, err := scanBinaryEntity()
	if err != nil {
		return nil, "", err
	}
	ent := gk.UploadEntityRequest{
		Name:        name.Name,
		Description: name.Description,
		Type:        model.BinaryData,
		Attr: model.MetaData{
			Binary: *attr,
		},
	}
	return &ent, fp, nil
}

func scanBinaryEntity() (*model.EntityTypeBinary, string, error) {
	fp := ""
	fmt.Print("Enter filepath: ")
	_, err := fmt.Fscan(os.Stdin, &fp)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(fp)
	if err != nil {
		return nil, "", err
	}
	bin := model.EntityTypeBinary{
		OriginalFilename: info.Name(),
		Size:             info.Size(),
		Extension:        filepath.Ext(fp),
	}
	return &bin, fp, nil
}

func (c *GkClient) addBinary() error {
	req, fp, err := newBinaryEntity()
	if err != nil {
		return err
	}
	guid, err := c.sendRequestBinary(req)
	if err != nil {
		return err
	}
	fmt.Println("Created entity DONE; GUID:", guid)
	err = c.sendBinaryData(fp, guid)
	if err != nil {
		return err
	}
	return nil
}

func (c *GkClient) sendRequestBinary(ent *gk.UploadEntityRequest) (string, error) {
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
