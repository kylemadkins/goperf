package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

// type Level string

// const (
// 	Info    Level = "INFO"
// 	Warning       = "WARNING"
// 	Error         = "ERROR"
// )

type Level byte

const (
	Info Level = iota + 1
	Warning
	Error
)

var (
	// 2023-04-12T12:12:40Z - WARNING - ...
	logRe      = regexp.MustCompile(`([^ ]+) - ([A-Z]+) - (.*)`)
	levelNames = map[string]Level{
		"INFO":    Info,
		"WARNING": Warning,
		"ERROR":   Error,
	}
)

type Log struct {
	Time    time.Time
	Level   Level
	Message string
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func parseLine(line string) (Log, error) {
	match := logRe.FindStringSubmatch(line)
	if match == nil {
		return Log{}, fmt.Errorf("bad log: %q", line)
	}

	t, err := parseTime(match[1])
	if err != nil {
		return Log{}, err
	}

	lvl, ok := levelNames[match[2]]
	if !ok {
		return Log{}, fmt.Errorf("invalid log level: %q", line)
	}

	log := Log{
		Time: t,
		// Level: Level(match[2]),
		Level:   lvl,
		Message: match[3],
	}

	return log, nil
}

// ParseLogs returns slices of Logs from r.
func ParseLogs(r io.Reader) ([]Log, error) {
	logs := make([]Log, 0, 50_000)
	// var logs []Log

	s := bufio.NewScanner(r)
	lNum := 0
	for s.Scan() {
		lNum++
		log, err := parseLine(s.Text())
		if err != nil {
			return nil, fmt.Errorf("%d: %w", lNum, err)
		}
		logs = append(logs, log)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
