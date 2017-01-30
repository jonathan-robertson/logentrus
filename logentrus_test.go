package logentrus_test

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/puddingfactory/logentrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)           // This will effect your stdout level, but not the level for LogentriesHook. You specify that priority on creation
	logrus.SetFormatter(&logrus.TextFormatter{}) // You an use any formatter; LogentriesHook will always format as JSON without interfering with your other hooks

	hook, err := logentrus.New(
		os.Getenv("LOGENTRIESTOKEN"), // fetching token from env vars here. You can make a token in your logentries account and are expected to have 1 token for each application
		&logentrus.Opts{ // include options (set to nil if options not necessary)
			Priority:        logrus.InfoLevel, // log level is inclusive. Setting to logrus.ErrorLevel, for example, would include errors, panics, and fatals, but not info or debug.
			TimestampFormat: "Jan 2 15:04:05", // timeFormat could be an empty string instead; doing so will default to logrus's typically time format.
			EncTLSConfig:    nil,              // setting config to nil means that conn will use root certs already set up on local system
			UnencryptedTCP:  true,             // disable encryption, but still use TCP
			UnencryptedUDP:  false,            // disable encryption and use UDP
			UnencryptedPort: 514,              // omitting will result in port 514 usage; valid options are 80, 514, and 10000
		},
	)
	if err != nil {
		panic(err)
	}
	logrus.AddHook(hook)
}

func TestDebug(t *testing.T) {
	logrus.Debug("This is a debug entry that should *not* show in logentries") // This won't appear in logentries due to the priority we set
}

func TestInfo(t *testing.T) {
	logrus.WithField("anotherField", "hi there!").Info("This is an info entry that should show up in logentries")
}

func TestError(t *testing.T) {
	logrus.WithField("the rent", "is too dang high").Error("This is an error entry that should also appear in logentries")
}
