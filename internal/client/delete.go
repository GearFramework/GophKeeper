package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func (c *GkClient) deleteEntity() error {
	var err error
	guid := ""
	fmt.Print("Enter GUID of deleting entity: ")
	_, err = fmt.Fscan(os.Stdin, &guid)
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"DELETE", c.Conf.Addr+"/v1/entities/"+guid, nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return NewError(resp.StatusCode)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Entity deleted success")
	return nil
}
