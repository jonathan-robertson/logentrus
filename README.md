# Logentrus | a [Logentries](https://logentries.com) hook for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/> [![GoDoc](https://godoc.org/github.com/puddingfactory/logentrus?status.svg)](https://godoc.org/github.com/puddingfactory/logentrus)

*Logrus created by [Simon Eskildsen](http://sirupsen.com)*

## Install

`go get -u github.com/sirupsen/logrus github.com/puddingfactory/logentrus`

## Setup

First, you should get a token for your logentries account, which you'll need to feed into your app somehow.

1. Log into your logentries account
1. Navigate to create Add New Log
1. Select Manual (specifically since we're using the Token-based approach)
1. Set your values such as your Log Name and Log Set
1. Receive your token!

I'd **strongly recommend against storing the token directly in your source code**. I personally use Environment Variables for testing purposes and have done so in the example provided below.

## Usage

Just like with Logrus, it's best to define your options in `init` or in some early stage of your program.

```go
package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/puddingfactory/logentrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{}) // logentrus hook will always submit JSON to Logentries

	hook, err := logentrus.New(
		os.Getenv("LOGENTRIESTOKEN"), // grabbing Logentries Token from environment variable
		&logentrus.Opts{
			Priority:        logrus.InfoLevel, // since set to InfoLevel, DebugLevel is the only level that will be ignored
			TimestampFormat: "Jan 2 15:04:05", // setting empty string here will default to logrus's typically time format
			EncTLSConfig:    nil,              // setting config to nil means that conn will use root certs from local system
			UnencryptedTCP:  false,            // disable encryption, but still use TCP
			UnencryptedUDP:  false,            // disable encryption and use UDP
			UnencryptedPort: 514,              // omitting will result in port 514 usage; valid options are 80, 514, and 10000
		},
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

Each field in logentrus.Opts is entirely optional and will have some default value if necessary.

Option | Description | Default | Valid options
--- | --- | --- | ---
Priority | set a threshold for log severity that should make it to Logentries | `logrus.DebugLevel` (all log types to be sent to Logentries) | `logrus.DebugLevel`, `logrus.InfoLevel`, `logrus.WarnLevel`, `logrus.ErrorLevel`, `logrus.FatalLevel`, `logrus.PanicLevel`
TimestampFormat | Change the timestamp format | logrus's default time format | `"Jan 2 15:04:05"`, or any format accepted by Golang
EncTLSConfig | provide a tls config to use embedded ca cert(s) | `nil` (use system's root certs) | see [this example](https://github.com/puddingfactory/logentrus#you-can-provide-your-own-set-of-root-certs-when-using-newencryptedhook)
UnencryptedTCP | `true` to disable encryption and still use TCP | `false` | `true` / `false`
UnencryptedUDP | `true` this to disable encryption and use UDP | `false` | `true` / `false`
UnencryptedPort | if using an unencrypted connection, choose a port here | `514` | `80`, `514`, and `10000`

Note that the entire logentrus.Opts param is optional as well. So if you're happy with the defaults, just enter `nil` as shown below:

```go
hook, err := logentrus.New(logentriesToken, nil)
```

## Features

### Logentrus does its own formatting

Since Logentries prefers JSON formatting, I didn't want to require it to be set in Logrus. Instead, there is a separate logrus.JSONFormatter within this hook that processes the log entries for you automatically.

You have the option of setting the logrus.JSONFormatter.TimestampFormatter value when calling logentrus.New if there's a Timestamp format you prefer.

```go
hook, err := logentrus.New(
	logentriesToken,
	&logentrus.Opts{
		TimestampFormat: "Jan 2 15:04:05",
	},
)
```

### Support for various transmission types

As a safety-net, data is sent with encrypted (TLS) over TCP by default.

But for those sitautions where you know that you'll only be transmitting non-sensitive data, then switching away from the default to an unencrypted TCP or UDP connection may have a noticable impact on the speed of your program - particularly if it's log-heavy.

```go
hook, err := logentrus.New(
	logentriesToken,
	&logentrus.Opts{
		UnencryptedTCP:  true,  // disable encryption, but still use TCP
		UnencryptedUDP:  false, // disable encryption and use UDP
		UnencryptedPort: 514,   // omitting will result in port 514 usage; valid options are 80, 514, and 10000
	},
)
```

### You can provide your own set of root certs when using NewEncryptedHook

This is a feature of Google's `crypto/tls` package that you can apply to logentrus [ref](https://golang.org/pkg/crypto/tls/#example_Dial).

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/puddingfactory/logentrus"
)

func main() {
	// Connecting with a custom root-certificate set.

	const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----`

	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	hook, err := logentrus.New(
		os.Getenv("LOGENTRIESTOKEN"),
		&logentrus.Opts{
			EncTLSConfig: &tls.Config{
				RootCAs: roots,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logrus.AddHook(hook)
}
```
