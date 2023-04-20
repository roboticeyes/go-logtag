package logtag

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type LogColor int

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

func ConfigureLogger(tags map[string]LogColor) {
	tagMap = tags
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) //remove timestamp, already included in grafana
}

// Error logs a message at level Error on the standard logger.
func addColoredTag(tag string, message string) string {
	col, ok := tagMap[tag]

	if !ok {
		return message
	}

	return toColoredText(col, "["+tag+"] ") + message
}

func toColoredText(col LogColor, message string) string {
	return col.ColorString() + message + Reset.ColorString()
}

func Printf(tag string, format string, v ...any) {
	log.Printf(addColoredTag(tag, format), v...)
}

func Println(tag string, msg string) {
	log.Print(addColoredTag(tag, msg))
}

func Errorf(tag string, format string, v ...any) {
	log.Printf(toColoredText(Red, "Error: ")+addColoredTag(tag, format), v)
}

func Error(tag string, msg string) {
	log.Print(toColoredText(Red, "Error: ") + addColoredTag(tag, msg))
}

func Warnf(tag string, format string, v ...any) {
	log.Printf(toColoredText(Yellow, "Warning: ")+addColoredTag(tag, format), v)
}

func Warn(tag string, msg string) {
	log.Print(toColoredText(Yellow, "Warning: ") + addColoredTag(tag, msg))
}

func Infof(tag string, format string, v ...any) {
	Printf(tag, format, v)
}

func Info(tag string, msg string) {
	Println(tag, msg)
}

func GinLogTag(tag string) gin.HandlerFunc {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := toColoredText(BrightBlue, c.Request.URL.Path)
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		dataLength := c.Writer.Size()

		if dataLength < 0 {
			dataLength = 0
		}

		method := toColoredText(BrightBlue, c.Request.Method)

		statusCodeString := fmt.Sprint(statusCode)
		if statusCode >= http.StatusInternalServerError {
			statusCodeString = toColoredText(Red, statusCodeString)
		} else if statusCode >= http.StatusBadRequest {
			statusCodeString = toColoredText(Yellow, statusCodeString)
		} else {
			statusCodeString = toColoredText(Green, statusCodeString)
		}

		if len(c.Errors) > 0 {
			Error(tag, c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s \"%s %s\" %s %d \"%s\" (%dms)", clientIP, hostname, method, path, statusCodeString, dataLength, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				Error(tag, msg)
			} else {
				Info(tag, msg)
			}
		}
	}
}
