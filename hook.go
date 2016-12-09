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
	AppName  string
	HostName string
	Priority logrus.Level

	Token string
	conn  net.Conn
}

const (
	host = "data.logentries.com"
	port = 443
)

// NewLogentriesHook creates and returns a new hook to an instance of logger.
// `hook, err := NewLogentriesHook("leServer", "leApp", "2bfbea1e-10c3-4419-bdad-7e6435882e1f", logrus.InfoLevel, nil)`
// `if err == nil { log.Hooks.Add(hook) }`
// NOTE: setting config to nil means that conn will use root certs already set up on local system
// Can provide own root certs by using example found here: https://golang.org/pkg/crypto/tls/#example_Dial
func NewLogentriesHook(hostName, appName, token string, priority logrus.Level, config *tls.Config) (*LogentriesHook, error) {
	tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	hook := &LogentriesHook{
		AppName:  appName,
		HostName: hostName,
		Priority: priority,
		Token:    token,
		conn:     tlsConn,
	}

	return hook, err
}

// Fire sends entry to Logentries
func (hook *LogentriesHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
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

// Levels returns the log levels supported by LogentriesHook
func (hook *LogentriesHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
