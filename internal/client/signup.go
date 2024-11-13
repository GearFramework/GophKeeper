package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func newSignupRequest() (*SignupRequest, error) {
	cred := SignupRequest{}
	fmt.Print("Enter username: ")
	_, err := fmt.Fscan(os.Stdin, &cred.Username)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Enter password fot %s: ", cred.Username)
	_, err = fmt.Fscan(os.Stdin, &cred.Password)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (c *GkClient) signup() error {
	cred, err := newSignupRequest()
	b, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"POST", c.Conf.Addr+"/v1/signup", bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return NewError(resp.StatusCode)
	}
	fmt.Println("Successfully user registration!")
	return nil
}
