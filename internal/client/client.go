package client

import (
	"fmt"
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
	token string
}

func (c *GkClient) auth() error {
	var err error
	c.Conf.Username, c.Conf.Password, err = scanSigninCredentials(c.Conf.Username)
	if err != nil {
		return err
	}
	c.token, err = c.Signin(c.Conf.Username, c.Conf.Password)
	if err != nil {
		return err
	}
	return nil
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
		return c.Signup()
	case CommandAdd:
		return c.Add()
	case CommandList:
		return c.ListEntities()
	case CommandView:
		return c.ViewEntity()
	case CommandUpload:
		return c.UploadEntity()
	case CommandDownload:
		return c.DownloadEntity()
	case CommandDel:
		return c.DeleteEntity()
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
