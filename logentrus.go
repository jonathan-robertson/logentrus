// Package logentrus acts as a Logentries (https://logentries.com) hook
// for Logrus (https://github.com/sirupsen/logrus)
package logentrus

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// Hook used to send logs to logentries
type Hook struct {
	Token    string
	Priority logrus.Level

	formatter *logrus.JSONFormatter
	conn      net.Conn
}

const (
	host = "data.logentries.com"
	port = 443
)

// New creates and returns a new hook for use in Logrus
func New(token, timestampFormat string, priority logrus.Level, config *tls.Config) (*Hook, error) {
	if token == "" {
		return nil, fmt.Errorf("Unable to create new LogentriesHook since a Token is required")
	}

	tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	hook := &Hook{
		Priority:  priority,
		Token:     token,
		formatter: &logrus.JSONFormatter{TimestampFormat: timestampFormat},
		conn:      tlsConn,
	}

	return hook, err
}

// Fire formats and sends JSON entry to Logentries service
func (hook *Hook) Fire(entry *logrus.Entry) error {
	line, err := hook.format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	if entry.Level <= hook.Priority {
		if _, err := hook.conn.Write([]byte(hook.Token + line)); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write to conn, %v", err)
			return err
		}
	}

	return nil
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
	return logrus.AllLevels
}
