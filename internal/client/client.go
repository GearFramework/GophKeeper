package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

// SigninRequest struct of auth user request
type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupRequest struct of register new user request
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GkClient strict of client
type GkClient struct {
	Conf  *Config
	Tls   *tls.Config
	token string
}

func (c *GkClient) auth() error {
	var err error
	c.Conf.Username, c.Conf.Password, err = scanSigninCredentials(c.Conf.Username)
	if err != nil {
		return err
	}
	c.token, err = c.signin(c.Conf.Username, c.Conf.Password)
	if err != nil {
		return err
	}
	return nil
}

func (c *GkClient) getTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: c.Tls,
	}
}

// Run start client
func (c *GkClient) Run() error {
	if c.Conf.Command.NeedAuth() {
		err := c.auth()
		if err != nil {
			return err
		}
	}
	switch c.Conf.Command {
	case CommandSignup:
		return c.signup()
	case CommandAdd:
		return c.add()
	case CommandList:
		return c.listEntities()
	case CommandView:
		return c.viewEntity()
	case CommandUpload:
		return c.uploadEntity()
	case CommandDownload:
		return c.downloadEntity()
	case CommandDel:
		return c.deleteEntity()
	}
	return nil
}

type nameEntity struct {
	Name        string
	Description string
}

func scanNameEntity() (*nameEntity, error) {
	name := nameEntity{}
	fmt.Print("Enter entity name: ")
	_, err := fmt.Fscan(os.Stdin, &name.Name)
	if err != nil {
		return nil, err
	}
	fmt.Print("Enter entity description: ")
	_, err = fmt.Fscan(os.Stdin, &name.Description)
	if err != nil {
		return nil, err
	}
	return &name, nil
}
