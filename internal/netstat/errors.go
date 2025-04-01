package netstat

import "errors"

var (
	errNetV6      = errors.New("v6 address - ignore")
	errParsingTCP = errors.New("error in parsing line TCP")
	errParsingUDP = errors.New("error in parsing line UDP")
)
