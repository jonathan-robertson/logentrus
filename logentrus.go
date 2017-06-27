// Package logentrus acts as a Logentries (https://logentries.com) hook
// for Logrus (https://github.com/Sirupsen/logrus)
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
	encrypt   bool
	token     string
	levels    []logrus.Level
	formatter *logrus.JSONFormatter
	network   string
	port      int
	tlsConfig *tls.Config
}

// Opts is a set of optional parameters for NewEncrytpedHook
type Opts struct {
	Priority        logrus.Level // defaults to logrus.DebugLevel (include all)
	TimestampFormat string       // defaults to logrus's typical timestamp format
	EncTLSConfig    *tls.Config  // defaults to use system's cert store; provide if you'd like to enforce your own root certs
	UnencryptedTCP  bool         // defaults to false (encryption enabled, using TCP)
	UnencryptedUDP  bool         // defaults to false (encryption enabled, using TCP)
	UnencryptedPort int          // defaults to 80; available ports are 80, 514, and 10000
}

const (
	version = "v2.0.0"
	host    = "data.logentries.com"
	tlsPort = 443
)

var (
	errTokenRequired = fmt.Errorf("unable to create new hook: a Token is required")
	errInvalidPort   = fmt.Errorf("unable to create new unencrypted hook: invalid port provided; only 80, 514, and 10000 are supported")
)

// New creates and returns a Logrus hook for Logentries Token-based logging
// ref: https://docs.logentries.com/docs/input-token
func New(token string, options *Opts) (hook *Hook, err error) {
	if token == "" {
		err = errTokenRequired
	} else {
		hook = &Hook{
			encrypt:   true,
			token:     token,
			levels:    logrus.AllLevels,
			formatter: &logrus.JSONFormatter{},
			network:   "tcp",
			port:      tlsPort,
		}

		if options != nil {
			hook.formatter.TimestampFormat = options.TimestampFormat
			hook.levels = logrus.AllLevels[:options.Priority+1]

			switch {
			case options.UnencryptedTCP:
				hook.encrypt = false
				hook.network = "tcp"
				hook.port = 514
			case options.UnencryptedUDP:
				hook.encrypt = false
				hook.network = "udp"
				hook.port = 514
			}

			if hook.encrypt {
				if options.EncTLSConfig != nil {
					hook.tlsConfig = options.EncTLSConfig
				}
			} else {
				switch options.UnencryptedPort {
				case 80, 514, 10000:
					hook.port = options.UnencryptedPort
				case 0: // ignore
				default:
					err = errInvalidPort
					return
				}
			}
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
		fmt.Fprintf(os.Stderr, "unable to read entry | err: %v | entry: %+v\n", err, entry)
		return err
	}

	if err = hook.write(line); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to write to conn | err: %v | line: %s\n", err, line)
	}

	return
}

// dial establishes a new connection which caller is responsible for closing
func (hook Hook) dial() (net.Conn, error) {
	if hook.encrypt {
		return tls.Dial(hook.network, fmt.Sprintf("%s:%d", host, hook.port), hook.tlsConfig)
	}
	return net.Dial(hook.network, fmt.Sprintf("%s:%d", host, hook.port))
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
