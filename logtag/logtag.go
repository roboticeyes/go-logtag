// Copyright (C) 2023 Robotic Eyes
//
// THIS CODE AND INFORMATION ARE PROVIDED "AS IS" WITHOUT WARRANTY OF ANY
// KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A
// PARTICULAR PURPOSE.

package logtag

import (
	"log"
	"time"
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
	Grey
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
	case Grey:
		return "\x1b[38;5;247m"
	case Reset:
		return "\x1b[0m"
	}

	return ""
}

var tagMap map[string]LogColor
var minLogLevel LogLevel = LevelInfo
var ignoreMap map[string]struct{}

func ConfigureLogger(tags map[string]LogColor, ignoreTags []string) {
	tagMap = tags
	ignoreMap = make(map[string]struct{})
	for _, tag := range ignoreTags {
		ignoreMap[tag] = struct{}{}
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) //remove timestamp, because we want to use colors
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

func addDateTime(message string) string {
	// add date/time using colors
	t := time.Now()
	timeString := t.Format("2006-01-02 15:04:05")

	return ToColoredText(BrightBlack, timeString) + message
}

func ToColoredText(col LogColor, message string) string {
	return col.ColorString() + message + Reset.ColorString()
}

func dontPrint(tag string, level LogLevel) bool {
	_, isIgnored := ignoreMap[tag]
	return isIgnored || minLogLevel > level
}

func Printf(tag string, format string, v ...any) {
	if dontPrint(tag, LevelInfo) {
		return
	}
	log.Printf(addDateTime(addColoredTag(tag, format)), v...)
}

func Println(tag string, msg string) {
	if dontPrint(tag, LevelInfo) {
		return
	}
	log.Print(addDateTime(addColoredTag(tag, msg)))
}

func Infof(tag string, format string, v ...any) {
	if dontPrint(tag, LevelInfo) {
		return
	}
	log.Printf(addColoredTag(tag, ToColoredText(Reset, "Info: ")+format), v...)
}

func Info(tag string, msg string) {
	if dontPrint(tag, LevelInfo) {
		return
	}
	log.Print(addColoredTag(tag, ToColoredText(Reset, "Info: ")+msg))
}

func Warnf(tag string, format string, v ...any) {
	if dontPrint(tag, LevelWarning) {
		return
	}
	log.Printf(addColoredTag(tag, ToColoredText(Yellow, "Warning: ")+format), v...)
}

func Warn(tag string, msg string) {
	if dontPrint(tag, LevelWarning) {
		return
	}
	log.Print(addColoredTag(tag, ToColoredText(Yellow, "Warning: ")+msg))
}

func Errorf(tag string, format string, v ...any) {
	if dontPrint(tag, LevelError) {
		return
	}
	log.Printf(addColoredTag(tag, ToColoredText(Red, "Error: ")+format), v...)
}

func Error(tag string, msg string) {
	if dontPrint(tag, LevelError) {
		return
	}
	log.Print(addColoredTag(tag, ToColoredText(Red, "Error: ")+msg))
}

func Fatalf(tag string, format string, v ...any) {
	if dontPrint(tag, LevelFatal) {
		return
	}
	log.Fatalf(addColoredTag(tag, ToColoredText(Red, "Fatal: ")+format), v...)
}

func Fatal(tag string, msg string) {
	if dontPrint(tag, LevelFatal) {
		return
	}
	log.Fatal(addColoredTag(tag, ToColoredText(Red, "Fatal: ")+msg))
}
