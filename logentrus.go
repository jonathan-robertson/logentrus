// Package logentrus acts as a Logentries (https://logentries.com) hook
// for Logrus (https://github.com/Sirupsen/logrus)
package logentrus

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/Sirupsen/logrus"
)

// Hook used to send logs to logentries
type Hook struct {
	levels    []logrus.Level
	token     string
	formatter *logrus.JSONFormatter
	tlsConfig *tls.Config
}

const (
	version = "v1.0.2"
	host    = "data.logentries.com"
	port    = 443
)

// New creates and returns a new hook for use in Logrus
func New(token, timestampFormat string, priority logrus.Level, config *tls.Config) (hook *Hook, err error) {
	if token == "" {
		err = fmt.Errorf("Unable to create new LogentriesHook since a Token is required")
	} else {
		hook = &Hook{
			levels:    logrus.AllLevels[:priority+1], // Include all levels at or within the provided priority
			token:     token,
			tlsConfig: config,
			formatter: &logrus.JSONFormatter{
				TimestampFormat: timestampFormat,
			},
		}

		// Test connection
		if conn, err := hook.dial(); err == nil {
			conn.Close()
		}

	}
	return
}

// Levels returns the log levels supported by this hook
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

// Fire formats and sends JSON entry to Logentries service
func (hook *Hook) Fire(entry *logrus.Entry) (err error) {
	line, err := hook.format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry | err: %v | entry: %+v\n", err, entry)
		return err
	}

	if err = hook.write(line); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to write to conn | err: %v | line: %s\n", err, line)
	}

	return
}

// dial establishes a new connection which caller responsible for closing
func (hook Hook) dial() (net.Conn, error) {
	return tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), hook.tlsConfig)
}

// write dials a connection and writes token and line in bytes to connection
func (hook *Hook) write(line string) (err error) {
	if conn, err := hook.dial(); err == nil {
		defer conn.Close()
		_, err = conn.Write([]byte(hook.token + line))
	}
	return
}

// format serializes the entry as JSON regardless of logentries's overall formatting
func (hook Hook) format(entry *logrus.Entry) (string, error) {
	serialized, err := hook.formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}
