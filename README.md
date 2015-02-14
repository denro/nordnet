# go-nordnet

Nordnet nEXT API Client. Comes as three separate packages:

## Installation

The APIClient is used for making REST requests, such as making orders.

  $ go get github.com/denro/nordnet/api

The Feed package is for reading realtime prices and trades.

  $ go get github.com/denro/nordnet/feed

Util contains authentication.

  $ go get github.com/denro/nordnet/util


## Usage

```go
package main

import (
	"github.com/denro/nordnet/api"
	"github.com/denro/nordnet/util"
)

var (
	pemData = []byte(`-----BEGIN PUBLIC KEY-----`) 
	user = []byte(`...`)
	pass = []byte(`...`)
)

func main() {
  cred, _ := util.GenerateCredentials(user, pass, pemData)
  client := api.NewAPIClient(*cred)
  client.Login()

  fmt.Println(client.Account())
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
