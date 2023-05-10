// Copyright (C) 2023 Robotic Eyes
//
// THIS CODE AND INFORMATION ARE PROVIDED "AS IS" WITHOUT WARRANTY OF ANY
// KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.

package logtag

import (
	"log"
)

type LogColor int
type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelWarning
	LevelError
	LevelFatal
)

const (
	Black LogColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
	Reset
)

func (c LogColor) ColorString() string {
	switch c {
	case Black:
		return "\x1b[38;30m"
	case Red:
		return "\x1b[38;31m"
	case Green:
		return "\x1b[38;32m"
	case Yellow:
		return "\x1b[38;33m"
	case Blue:
		return "\x1b[38;34m"
	case Magenta:
		return "\x1b[38;35m"
	case Cyan:
		return "\x1b[38;36m"
	case White:
		return "\x1b[38;37m"
	case BrightBlack:
		return "\x1b[38;90m"
	case BrightRed:
		return "\x1b[38;91m"
	case BrightGreen:
		return "\x1b[38;92m"
	case BrightYellow:
		return "\x1b[38;93m"
	case BrightBlue:
		return "\x1b[38;94m"
	case BrightMagenta:
		return "\x1b[38;95m"
	case BrightCyan:
		return "\x1b[38;96m"
	case BrightWhite:
		return "\x1b[38;97m"
	case Reset:
		return "\x1b[0m"
	}

	return ""
}

var tagMap map[string]LogColor
var minLogLevel LogLevel = LevelInfo

func ConfigureLogger(tags map[string]LogColor) {
	tagMap = tags
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) //remove timestamp, already included in grafana
}

func SetMinimumLogLevel(l LogLevel) {
	minLogLevel = l
}

// Error logs a message at level Error on the standard logger.
func addColoredTag(tag string, message string) string {
	col, ok := tagMap[tag]

	if !ok {
		return message
	}

	return ToColoredText(col, "["+tag+"] ") + message
}

func ToColoredText(col LogColor, message string) string {
	return col.ColorString() + message + Reset.ColorString()
}

func Printf(tag string, format string, v ...any) {
	if minLogLevel > LevelInfo {
		return
	}
	log.Printf(addColoredTag(tag, format), v...)
}

func Println(tag string, msg string) {
	if minLogLevel > LevelInfo {
		return
	}
	log.Print(addColoredTag(tag, msg))
}

func Infof(tag string, format string, v ...any) {
	if minLogLevel > LevelInfo {
		return
	}
	Printf(tag, format, v...)
}

func Info(tag string, msg string) {
	if minLogLevel > LevelInfo {
		return
	}
	Println(tag, msg)
}

func Warnf(tag string, format string, v ...any) {
	if minLogLevel > LevelWarning {
		return
	}
	log.Printf(addColoredTag(tag, ToColoredText(Yellow, "Warning: ")+format), v...)
}

func Warn(tag string, msg string) {
	if minLogLevel > LevelWarning {
		return
	}
	log.Print(addColoredTag(tag, ToColoredText(Yellow, "Warning: ")+msg))
}

func Errorf(tag string, format string, v ...any) {
	if minLogLevel > LevelError {
		return
	}
	log.Printf(addColoredTag(tag, ToColoredText(Red, "Error: ")+format), v...)
}

func Error(tag string, msg string) {
	if minLogLevel > LevelError {
		return
	}
	log.Print(addColoredTag(tag, ToColoredText(Red, "Error: ")+msg))
}

func Fatalf(tag string, format string, v ...any) {
	if minLogLevel > LevelFatal {
		return
	}
	log.Fatalf(addColoredTag(tag, ToColoredText(Red, "Fatal: ")+format), v...)
}

func Fatal(tag string, msg string) {
	if minLogLevel > LevelFatal {
		return
	}
	log.Fatal(addColoredTag(tag, ToColoredText(Red, "Fatal: ")+msg))
}
