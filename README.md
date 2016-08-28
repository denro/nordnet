# Nordnet

[![Build Status](https://travis-ci.org/denro/nordnet.svg?branch=master)](https://travis-ci.org/denro/nordnet)
[![GoDoc](https://godoc.org/github.com/denro/nordnet?status.svg)](http://godoc.org/github.com/denro/nordnet)
[![Go Report Card](https://goreportcard.com/badge/github.com/denro/nordnet)](https://goreportcard.com/report/github.com/denro/nordnet)

Go implementation of the Nordnet External API.

https://api.test.nordnet.se/api-docs/index.html


## Installation

`go get github.com/denro/nordnet`

## Usage


### REST API Client

```go
package main

import (
	"fmt"
	"github.com/denro/nordnet/api"
	"github.com/denro/nordnet/util"
)

var (
	pemData = []byte(`-----BEGIN PUBLIC KEY-----`)
	user    = []byte(`...`)
	pass    = []byte(`...`)
)

func main() {
	cred, _ := util.GenerateCredentials(user, pass, pemData)
	client := api.NewAPIClient(cred)
	client.Login()

	fmt.Println(client.Accounts())
}
```

### Feed Client

```go
package main

import (
	"fmt"
	"github.com/denro/nordnet/feed"
)

var (
	sessionKey = "..."
	address    = "..."
)

func main() {
	feed, _ := feed.NewPrivateFeed(address)
	feed.Login(sessionKey, nil)

	msgChan := make(chan *PrivateMsg)
	errChan := make(chan error)
	feed.Dispatch(msgChan, errChan)

	for _, msg := range msgChan {
		fmt.Println(msg)
	}
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
