package config

import "flag"

var (
	Port = flag.Int("port", 8080, "The server port")
)

func Start() error {
	flag.Parse()

	return nil
}
