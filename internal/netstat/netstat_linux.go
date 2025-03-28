//go:build linux

package netstat

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

const (
	lineDelim      = "\n"
	netstatCommand = "netstat"
)

var (
	regTCP      = regexp.MustCompile(`tcp.+:(\d+)\s.*:(\*|\d+)\s+([[:graph:]]+)\s+([\w-_@\.]+)\s+\d+\s+([[:graph:]]*).*`)
	regUDP      = regexp.MustCompile(`udp.+:(\d+)\s.*:(\*|\d+)\s+([\w-_@\.]+)\s+\d+\s+([[:graph:]]*).*`)
	netstatARGS = []string{"-apeW", "-A", "inet", "--numeric-hosts", "--numeric-ports"}
)

func parseLineTCP(line string) (Socket, string, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	splitLine := regTCP.FindStringSubmatch(line)

	if len(splitLine) != 6 { // столько должно быть распарсенных элементов
		return sock, "", errParsingTCP
	}
	slog.Debug("Netstat", "port", splitLine[1], "status", splitLine[3], "user", splitLine[4], "pid", splitLine[5])

	sock.User = splitLine[4]

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32) //nolint:gosec
	}

	if splitLine[5] != "-" {
		s2 := strings.Split(splitLine[5], "/")
		if len(s2) != 2 { // Не в формате
			return sock, splitLine[3], nil // Значения по умолчанию
		}
		if i32, err := strconv.Atoi(s2[0]); err == nil {
			sock.Pid = int32(i32) //nolint:gosec
		}
		sock.Command = s2[1]
	}

	return sock, splitLine[3], nil
}

func parseLineUDP(line string) (Socket, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	splitLine := regUDP.FindStringSubmatch(line)

	if len(splitLine) != 5 { // столько должно быть распарсенных элементов
		return sock, errParsingUDP
	}
	slog.Debug("Netstat", "port", splitLine[1], "user", splitLine[3], "pid", splitLine[4])

	sock.User = splitLine[3]

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32) //nolint:gosec
	}

	if splitLine[4] != "-" {
		s2 := strings.Split(splitLine[4], "/")
		if len(s2) != 2 { // Не в формате
			return sock, nil // Значения по умолчанию
		}
		if i32, err := strconv.Atoi(s2[0]); err == nil {
			sock.Pid = int32(i32) //nolint:gosec
		}
		sock.Command = s2[1]
	}

	return sock, nil
}
