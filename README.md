# Logrusentries | a helpful hook for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>
Logrusentries is a [Logentries](https://logentries.com) hook for [Logrus](https://github.com/sirupsen/logrus).

*Logrus created by [Simon Eskildsen](http://sirupsen.com)*

# Install
`go get -u github.com/sirupsen/logrus github.com/puddingfactory/logrus-logentries-hook`

# Usage
First, you should get a token for your logentries account, which you'll need to feed into your app somehow.

1. Log into your logentries account
2. Navigate to create Add New Log
3. Select Manual (specifically since we're using TCP directly)
4. Set your values such as your Log Name and Log Set
5. Receive your token!

I'd **strongly recommend against storing the token directly in your source code**.
I personally use Environment Variables for testing purposes and have done so in the example provided below.

Just like with Logrus, it's best to define your options and attach this hook within `init` or in some early stage of your program.

```go
package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/puddingfactory/logrusentries"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{}) // hook will always format as JSON with its own formatter

	hook, err := logrusentries.New(
		os.Getenv("TOKEN"), // grabbing this from environment variable
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
## Hook does its own formatting
tl;dr: Hook has its own logrus.JSONFormatter that's automatically set and called exactly the same way that it is for Logrus.
You can set the TimestampFormat when calling NewLogentriesHook.

#### Why?
The TextLogger option for Logrus is incredibly well formatted and I am a huge fan of it. It even colors the prefix (grey for `DEBU`, blue for `INFO`, red for `ERR`, etc)

```
DEBU[0001] This is another debug
INFO[0001] This is another info
ERRO[0001] This is another error
```

Unfortunately, sending data formatted with TextLogger to Logentries is... ugly:

```
[37mDEBU[0m[0000] This is debug
[34mINFO[0m[0000] This is info
[31mERRO[0m[0000] This is an error
```

Logentries has several display formats available for JSON-formatted data, so sending the entry as JSON string is a no-brainer.

