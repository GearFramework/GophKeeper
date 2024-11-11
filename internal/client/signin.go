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
)

func scanSigninCredentials(username string) (string, string, error) {
	if username == "\x00" {
		fmt.Print("Enter username: ")
		_, err := fmt.Fscan(os.Stdin, &username)
		if err != nil {
			return "", "", err
		}
	}
	fmt.Print("Enter password: ")
	password := ""
	if _, err := fmt.Fscan(os.Stdin, &password); err != nil {
		return "", "", err
	}
	return username, password, nil
}

// Signin authority in remote server
func (c *GkClient) Signin(username, password string) (string, error) {
	b, err := json.Marshal(SigninRequest{username, password})
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"POST", c.Conf.Addr+"/v1/signin", bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
