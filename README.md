go-nordnet
==========

Nordnet nEXT API Client

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
