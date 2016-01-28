// Simple example which authenticate, a few rest api requests and then starts a subscription to trades and depths
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"

	"github.com/denro/nordnet/api"
	"github.com/denro/nordnet/feed"
	"github.com/denro/nordnet/util"
)

func main() {

	pretty := func (v interface{}) string {
		b, _ := json.Marshal(v)
		return fmt.Sprintf("%s\n",string(b))
	}

	// Never hardcode secrets to your code
	// export NORDNET_USER="..."
	user := []byte(os.Getenv("NORDNET_USER"))
	// export NORDNET_PASS="..."
	pass := []byte(os.Getenv("NORDNET_PASS"))
	// export NORDNET_PEMDATA=`cat /where/your/pem/file/is.pem`
	pemData := []byte(os.Getenv("NORDNET_PEMDATA"))

	cred, err := util.GenerateCredentials(user, pass, pemData)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := api.NewAPIClient(cred)
	session, err := client.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Take the rest api for a ride
	status, _ := client.SystemStatus()
	pretty(status)

	accounts, _ := client.Accounts()
	for _, account := range accounts {
		accountinfo, _ := client.Account(account.Accno)
		pretty(accountinfo)
	}

	instrument, _ := client.SearchInstruments(&api.Params{
    "query": "volvo-b",
	})
	pretty(instrument)

	// Open private feed
	privAddr := fmt.Sprintf("%s:%d", session.PrivateFeed.Hostname, session.PrivateFeed.Port)
	privfeed, _ := feed.NewPrivateFeed(privAddr)
	privfeed.Login(client.SessionKey, nil)

	// Open public feed
	pubAddr := fmt.Sprintf("%s:%d", session.PublicFeed.Hostname, session.PublicFeed.Port)
	pubfeed, _ := feed.NewPublicFeed(pubAddr)
	pubfeed.Login(client.SessionKey, nil)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	privmsgChan := make(chan *feed.PrivateMsg, 1000)
	pubmsgChan := make(chan *feed.PublicMsg, 1000)
	errChan := make(chan error, 1000)
	privfeed.Dispatch(privmsgChan, errChan)
	pubfeed.Dispatch(pubmsgChan, errChan)

	// Start subscriptions
	pubfeed.Subscribe(feed.PriceArgs{T: "price", I: "101", M: 11})
	pubfeed.Subscribe(feed.DepthArgs{T: "depth", I: "101", M: 11})

	// Receive messages until exit channel is messaged
	for {
		select {
		case msg := <-privmsgChan:
			fmt.Printf("Private feed: %s\n", pretty(msg))
		case msg := <-pubmsgChan:
			fmt.Printf("Public feed: %s\n", pretty(msg))
		case msg := <-errChan:
			fmt.Printf("Error chan: %s\n", pretty(msg))
		case <-exit:
			privfeed.Close()
			pubfeed.Close()
			client.Logout()
			return
		}
	}
}
