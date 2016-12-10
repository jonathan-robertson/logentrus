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
		os.Getenv("TOKEN"), // fetching token from env vars here. You can make a token in your logentries account and are expected to have 1 token for each application
		"Jan 2 15:04:05",   // timeFormat could be an empty string instead; doing so will default to logrus's typically time format.
		logrus.InfoLevel,   // log level is inclusive. Setting to logrus.ErrorLevel, for example, would include errors, panics, and fatals, but not info or debug.
		nil,                // setting config to nil means that conn will use root certs already set up on local system
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
