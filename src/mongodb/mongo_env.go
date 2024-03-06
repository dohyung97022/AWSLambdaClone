package mongodb

import "os"

type mongoEnv struct {
	username string
	password string
	dns      string
	port     string
}

var env mongoEnv

func init() {
	env = mongoEnv{
		username: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		password: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		dns:      os.Getenv("MONGO_DNS"),
		port:     os.Getenv("MONGO_PORT"),
	}
}
