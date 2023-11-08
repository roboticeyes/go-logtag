package logtag_gin

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roboticeyes/go-logtag/logtag"
)

type MethodAndPath struct {
	HttpMethod string
	Path       string
}

func GinLogTag(tag string, ignorePaths []MethodAndPath) gin.HandlerFunc {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path

		if c.Request.URL.Query().Encode() != "" {
			path = c.Request.URL.Path + "?" + c.Request.URL.Query().Encode()
		}
		coloredPath := logtag.ToColoredText(logtag.BrightBlue, path)
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

		method := logtag.ToColoredText(logtag.BrightBlue, c.Request.Method)

		statusCodeString := fmt.Sprint(statusCode)
		if statusCode > http.StatusInternalServerError {
			statusCodeString = logtag.ToColoredText(logtag.Red, statusCodeString)
		} else if statusCode > http.StatusBadRequest {
			statusCodeString = logtag.ToColoredText(logtag.Yellow, statusCodeString)
		} else {
			statusCodeString = logtag.ToColoredText(logtag.Green, statusCodeString)
		}

		if len(c.Errors) > 0 {
			logtag.Error(tag, c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if contains(ignorePaths, c.Request.Method, path) && statusCode < 300 {
				return
			}

			msg := fmt.Sprintf("\"%s %s\" code=%s %d \"%s\" (%dms) %s - %s ", method, coloredPath, statusCodeString, dataLength, clientUserAgent, latency, clientIP, hostname)
			if statusCode >= http.StatusInternalServerError {
				logtag.Error(tag, msg)
			} else {
				logtag.Info(tag, msg)
			}
		}
	}
}

func contains(list []MethodAndPath, method, path string) bool {
	for _, entry := range list {
		if entry.HttpMethod == method {
			matched, err := regexp.Match(entry.Path, []byte(path))
			if matched {
				return true
			}
			if err != nil {
				logtag.Error("ERRRROR", err.Error())
			}
		}
	}
	return false
}
