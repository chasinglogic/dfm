package logger

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LevelWarn LogLevel = iota
	LevelInfo
	LevelVerbose
	LevelDebug
)

func getActiveLevel() LogLevel {
	switch os.Getenv("DFM_LOG_LEVEL") {
	case "DEBUG":
		return LevelDebug
	case "WARN":
		return LevelWarn
	case "VERBOSE":
		return LevelVerbose
	default:
		return LevelInfo
	}
}

var activeLevel = getActiveLevel()

type LevelLogger struct {
	level  LogLevel
	stream io.Writer
}

func (ll *LevelLogger) Write(p []byte) (int, error) {
	if activeLevel >= ll.level {
		return ll.stream.Write(p)
	}

	return 0, nil
}

var Debug = log.New(
	&LevelLogger{
		level:  LevelDebug,
		stream: os.Stderr,
	},
	"DEBUG ",
	log.Lshortfile,
)

var Verbose = log.New(
	&LevelLogger{
		level:  LevelVerbose,
		stream: os.Stderr,
	},
	"",
	0,
)

var Info = log.New(
	&LevelLogger{
		level:  LevelInfo,
		stream: os.Stderr,
	},
	"INFO ",
	0,
)

var Warn = log.New(
	&LevelLogger{
		level:  LevelWarn,
		stream: os.Stderr,
	},
	"WARNING ",
	0,
)
