//go:build windows

package netstat

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

const (
	lineDelim      = "\r\n"
	netstatCommand = `C:\Windows\System32\netstat.exe`
)

var (
	regTCP      = regexp.MustCompile(`.*TCP.+:(\d+)\s.*:(\*|\d+)\s+([[:graph:]]+)\s+(\d+).*`)
	regUDP      = regexp.MustCompile(`.*UDP.+:(\d+)\s.*:(\*|\d+)\s+(\d+).*`)
	netstatARGS = []string{"-ano"}
)

func parseLineTCP(line string) (Socket, string, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	if strings.Contains(line, "[") {
		return sock, "", errNetV6
	}

	splitLine := regTCP.FindStringSubmatch(line)

	if len(splitLine) != 5 { // столько должно быть распарсенных элементов
		return sock, "", errParsingTCP
	}
	slog.Debug("Netstat", "port", splitLine[1], "status", splitLine[3], "pid", splitLine[4])

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32) //nolint:gosec
	}

	if i32, err := strconv.Atoi(splitLine[4]); err == nil {
		sock.Pid = int32(i32) //nolint:gosec
	}

	return sock, splitLine[3], nil
}

func parseLineUDP(line string) (Socket, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	if strings.Contains(line, "[") {
		return sock, errNetV6
	}

	splitLine := regUDP.FindStringSubmatch(line)

	if len(splitLine) != 4 { // столько должно быть распарсенных элементов
		return sock, errParsingUDP
	}
	slog.Debug("Netstat", "port", splitLine[1], "pid", splitLine[3])

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32) //nolint:gosec
	}

	if i32, err := strconv.Atoi(splitLine[3]); err == nil {
		sock.Pid = int32(i32) //nolint:gosec
	}

	return sock, nil
}
