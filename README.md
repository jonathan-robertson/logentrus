# Logentrus | a [Logentries](https://logentries.com) hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/> [![GoDoc](https://godoc.org/github.com/puddingfactory/logentrus?status.svg)](https://godoc.org/github.com/puddingfactory/logentrus)
*Logrus created by [Simon Eskildsen](http://sirupsen.com)*

# Install
`go get -u github.com/Sirupsen/logrus github.com/puddingfactory/logentrus`

# Setup
First, you should get a token for your logentries account, which you'll need to feed into your app somehow.

1. Log into your logentries account
2. Navigate to create Add New Log
3. Select Manual (specifically since we're using TCP directly)
4. Set your values such as your Log Name and Log Set
5. Receive your token!

I'd **strongly recommend against storing the token directly in your source code**. I personally use Environment Variables for testing purposes and have done so in the example provided below.

# Usage
Just like with Logrus, it's best to define your options in `init` or in some early stage of your program.

```go
package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/puddingfactory/logentrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{}) // logentrus hook will always submit JSON to Logentries

	hook, err := logentrus.New(
		os.Getenv("TOKEN"), // grabbing Logentries Token from environment variable
		"Jan 2 15:04:05",   // setting empty string here will default to logrus's typically time format
		logrus.InfoLevel,   // since set to InfoLevel, DebugLevel is the only level that will be ignored
		nil,                // setting config to nil means that conn will use root certs from local system
	)
	if err != nil {
		panic(err)
	}
	logrus.AddHook(hook)
}

func main() {
	logrus.Debug("This is a debug entry that should *not* show in logentries")
	logrus.Info("This is an info entry that should show up in logentries")
}
```

# Features
## Logentrus does its own formatting
Since Logentries prefers JSON formatting, I didn't want to require it to be set in Logrus. Instead, there is a separate logrus.JSONFormatter within this hook that processes the log entries for you automatically.

You have the option of setting the logrus.JSONFormatter.TimestampFormatter value when calling logentrus.New if there's a Timestamp format you prefer.

## You can provide your own set of root certs
This is a feature of Google's `crypto/tls` package that you can apply to logentrus.

1. First you'll want to create a `tls.Config` by following this example: https://golang.org/pkg/crypto/tls/#example_Dial
2. After that, you can drop your `tls.Config` into logentrus.New.
