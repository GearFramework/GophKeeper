package db

// ConnectionConfig config for connection to db
type ConnectionConfig struct {
	ConnectionDSN   string
	ConnectMaxOpens int
}
