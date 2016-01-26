package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/denro/nordnet/api"
	"github.com/denro/nordnet/feed"
	"github.com/denro/nordnet/util"
)

func main() {
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
	fmt.Printf("\n%#v\n\n", status)

	accounts, _ := client.Accounts()
	for _, account := range accounts {
		accountinfo, _ := client.Account(account.Accno)
		fmt.Printf("%#v\n\n", accountinfo)
	}

	markets, _ := client.Markets()
	fmt.Printf("%#v\n\n", markets)

	// Open private feed
	privAddr := fmt.Sprintf("%s:%d", session.PrivateFeed.Hostname, session.PrivateFeed.Port)
	privfeed, _ := feed.NewPrivateFeed(privAddr)
	privfeed.Login(client.SessionKey, nil)

	// Open public feed
	pubAddr := fmt.Sprintf("%s:%d", session.PublicFeed.Hostname, session.PublicFeed.Port)
	pubfeed, _ := feed.NewPublicFeed(pubAddr)
	pubfeed.Login(client.SessionKey, nil)

	privmsgChan := make(chan *feed.PrivateMsg)
	pubmsgChan := make(chan *feed.PublicMsg)
	errChan := make(chan error)

	privfeed.Dispatch(privmsgChan, errChan)
	pubfeed.Dispatch(pubmsgChan, errChan)

	pubfeed.Subscribe(&feed.TradingStatusArgs{T: "trading_status", I: "1869", M: 30})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-privmsgChan:
			fmt.Printf("Private message: %+v\n", msg)
		case msg := <-pubmsgChan:
			fmt.Printf("Public message: %+v\n", msg)
		case msg := <-errChan:
			fmt.Printf("Error message: %+v\n", msg)
		case <-exit:
			privfeed.Close()
			pubfeed.Close()
			client.Logout()
			return
		}
	}
}
