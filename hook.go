package logrus_logentries

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/Sirupsen/logrus"
)

// LogentriesHook to send logs via logentries service
type LogentriesHook struct {
	Token    string
	Priority logrus.Level

	formatter *logrus.JSONFormatter
	conn      net.Conn
}

const (
	host = "data.logentries.com"
	port = 443
)

// NewLogentriesHook creates and returns a new hook to an instance of logger.
// `hook, err := NewLogentriesHook("2bfbea1e-10c3-4419-bdad-7e6435882e1f", "Jan 2 15:04:05", logrus.InfoLevel, nil)`
// `if err == nil { log.Hooks.Add(hook) }`
// Can provide own root certs by using example found here: https://golang.org/pkg/crypto/tls/#example_Dial
func NewLogentriesHook(token, timestampFormat string, priority logrus.Level, config *tls.Config) (*LogentriesHook, error) {
	if token == "" {
		return nil, fmt.Errorf("Unable to create new LogentriesHook since a Token is required")
	}

	tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	hook := &LogentriesHook{
		Priority:  priority,
		Token:     token,
		formatter: &logrus.JSONFormatter{TimestampFormat: timestampFormat},
		conn:      tlsConn,
	}

	return hook, err
}

// Fire sends entry to Logentries
func (hook *LogentriesHook) Fire(entry *logrus.Entry) error {
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

func (hook LogentriesHook) format(entry *logrus.Entry) (string, error) {
	serialized, err := hook.formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// Levels returns the log levels supported by LogentriesHook
func (hook *LogentriesHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
