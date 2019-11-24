package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/LiangXianSen/gin-demo/errors"
)

type F logrus.Fields

type Logger struct {
	pipeline F
	*logrus.Logger
}

// Write writes content info logger Fields
func (log *Logger) Write(msg string, fd interface{}) {
	log.pipeline[msg] = fd
}

// NewLogger returns logrus instance
func NewLogger() *Logger {
	return &Logger{
		pipeline: F{},
		Logger:   logrus.New(),
	}
}

type LoggerOptions struct {
	Application  string
	Version      string
	EnableOutput bool
	EnableDebug  bool
}

// Logger is a middleware which provide a logger in ctx.
// Records each handling on os.stdout.
// nolint:funlen
func LoggerM(opt LoggerOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		logS := NewLogger()
		logS.SetFormatter(&logrus.JSONFormatter{})
		logS.SetOutput(os.Stdout)
		logS.SetLevel(logrus.InfoLevel)

		if opt.EnableDebug {
			logS.SetLevel(logrus.DebugLevel)
		}
		c.Set("logger", logS)

		// Replace gin writer for backup writer stream
		writer := new(multiWriter)
		writer.ctx = c
		writer.ResponseWriter = c.Writer
		c.Writer = writer

		c.Next()

		statusCode := c.Writer.Status()
		requestID := c.GetString("request_id")
		duration := Milliseconds(time.Since(start))

		info := F{
			"start":       start,
			"path":        path,
			"method":      method,
			"client_ip":   clientIP,
			"status_code": statusCode,
			"request_id":  requestID,
			"runtime":     duration,
			"version":     opt.Version,
			"application": opt.Application,
		}

		resp, _ := c.Get("response")
		if body, ok := resp.(map[string]interface{}); ok {
			info["response"] = F(body)
		} else {
			info["response"] = resp
		}

		// Writes data
		for k, v := range logS.pipeline {
			if _, ok := info[k]; !ok {
				info[k] = v
			}
		}

		filterBodyTooLong(info)

		if err, ok := c.Get("error"); ok {
			info["error"] = fmt.Sprintf("%v", err)
			if opt.EnableDebug {
				if e, ok := err.(*errors.Error); ok && e.Stack() != nil {
					info["error"] = fmt.Sprintf("%+v", err.(*errors.Error).Stack())
				}
			}
			logS.WithFields(logrus.Fields(info)).Error("error occurred")
			return
		}

		if opt.EnableOutput {
			logS.WithFields(logrus.Fields(info)).Info("finished")
		}
	}
}

// multiWriter is a backup of gin responseWriter
type multiWriter struct {
	gin.ResponseWriter
	ctx *gin.Context
}

func (w *multiWriter) isJSONResponse() bool {
	return strings.Contains(w.Header().Get("Content-Type"), "application/json")
}

// Write writes response then backups to ctx
func (w *multiWriter) Write(b []byte) (int, error) {
	var resp F
	if w.isJSONResponse() {
		if err := json.Unmarshal(b, &resp); err != nil {
			return 0, err
		}
	}
	if len(resp) == 0 {
		resp = F{"body": string(b)}
	}
	w.ctx.Set("response", map[string]interface{}(resp))
	return w.ResponseWriter.Write(b)
}

func Milliseconds(t time.Duration) float64 {
	m := t.Seconds() * 1000
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", m), 64)
	return f
}

const maxLengthToFilter = 512

func filterBodyTooLong(fields F) {
	for k, v := range fields {
		switch obj := v.(type) {
		case string:
			if len(obj) > maxLengthToFilter {
				fields[k] = "SIZE(" + strconv.Itoa(len(obj)) + ")"
			}
		case F:
			filterBodyTooLong(obj)
		case []byte:
			fields[k] = "SIZE(" + strconv.Itoa(len(obj)) + ")"
		default:
			// []interface{}
			// Slice类型,并且每个元素类型必须一致
			s := reflect.ValueOf(v)
			if v != nil && s.Kind() == reflect.Slice {
				for i := 0; i < s.Len(); i++ {
					if mp, ok := s.Index(i).Interface().(F); ok {
						filterBodyTooLong(mp)
					}
					if mp, ok := s.Index(i).Interface().(string); ok {
						if len(mp) > maxLengthToFilter {
							if p, ok := s.Index(i).Addr().Interface().(*string); ok {
								*p = "SIZE(" + strconv.Itoa(len(mp)) + ")"
							}
							if p, ok := s.Index(i).Addr().Interface().(*interface{}); ok {
								*p = "SIZE(" + strconv.Itoa(len(mp)) + ")"
							}
						}
					}
				}
			}
		}
	}
}
