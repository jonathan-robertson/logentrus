// Package logentrus acts as a Logentries (https://logentries.com) hook
// for Logrus (https://github.com/Sirupsen/logrus)
package logentrus

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

// Hook used to send logs to logentries
type Hook struct {
	levels    []logrus.Level
	token     string
	formatter *logrus.JSONFormatter
	tlsConfig *tls.Config
	conn      net.Conn
}

const (
	host       = "data.logentries.com"
	port       = 443
	retryCount = 3
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

		err = hook.dial()
	}
	return
}

// Fire formats and sends JSON entry to Logentries service
func (hook *Hook) Fire(entry *logrus.Entry) (err error) {
	line, err := hook.format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry | err: %v | entry: %+v\n", err, entry)
		return err
	}

	if err = hook.write(line); err != nil { // First write attempt
		for i := 0; i < retryCount; i++ {
			time.Sleep(time.Second) // Rest 1 second between retries
			fmt.Fprintf(os.Stderr, "WARNING: Trouble writing to conn; retrying %d of %d | err: %v\n", i, retryCount, err)
			if dialErr := hook.dial(); dialErr != nil { // Problem with write, so dial new connection and retry if possible
				fmt.Fprintf(os.Stderr, "ERROR: Unable to dial new connection | dialErr: %v\n", dialErr)
				return err
			}
			if err = hook.write(line); err == nil { // Retry write
				break
			}
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to write to conn | err: %v | line: %s\n", err, line)
	}

	return
}

func (hook *Hook) dial() (err error) {
	hook.conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), hook.tlsConfig)
	return
}

func (hook Hook) write(line string) (err error) {
	_, err = hook.conn.Write([]byte(hook.token + line))
	return
}

func (hook Hook) format(entry *logrus.Entry) (string, error) {
	serialized, err := hook.formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// Levels returns the log levels supported by this hook
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}
