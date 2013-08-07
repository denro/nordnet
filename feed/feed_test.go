package feed

import (
	"github.com/denro/go-nordnet/api"
	"testing"
)

var client = setupClient()

func TestNewFeed(t *testing.T) {
	_, err := NewFeed(address, service, sessionKey)
	if err != nil {
		t.Fatal(err)
	}
}

func setupClient() *api.APIClient {
	client := &api.APIClient{}
	client.URL = "https://api.test.nordnet.se/next"
	client.SessionKey = "test"

	return client
}
