# go-nordnet

Nordnet nEXT API Client. Comes as three separate packages:

## Installation

The APIClient is used for making REST requests, such as making orders.

  $ go get github.com/denro/go-nordnet/api

The Feed package is for reading realtime prices and trades.

  $ go get github.com/denro/go-nordnet/feed

Util contains authentication.

  $ go get github.com/denro/go-nordnet/util


## Usage

```go
package main

import (
	"github.com/denro/go-nordnet/api"
	"github.com/denro/go-nordnet/util"
	"log"
)

var (
	pemData = []byte(`-----BEGIN PUBLIC KEY-----`) 
	user = []byte(`...`)
	pass = []byte(`...`)
)

func main() {
  cred, _ := util.GenerateCredentials(user, pass, pemData)
  client := api.NewAPIClientwLogin(*cred)
  log.Println(client.Accounts())
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
