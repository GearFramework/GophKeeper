package client

// Config of client cli-application
type Config struct {
	Command  Command
	Addr     string
	Username string
	Password string
	Type     string
}
