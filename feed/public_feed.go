package feed

import (
	"encoding/json"
	"io"
	"log"
)

type PublicFeed struct {
	*Feed
}

type FeedArgs struct {
	T string `json:"t"`
	M int64  `json:"m"`
	I string `json:"i"`
}

func NewPublicFeed(address, service, sessionKey string) (*PublicFeed, error) {
	f, err := newFeed(address, service, sessionKey)
	if err != nil {
		return nil, err
	}

	pubf := &PublicFeed{Feed: f}

	if err = pubf.Login(); err != nil {
		return nil, err
	}

	return pubf, nil
}

func (f *PublicFeed) DispatchListener(feedChan chan *FeedCmd) {
	go listenOn(f.Conn, feedChan)
}

func (f *PublicFeed) Subscribe(args interface{}) error {
	return f.WriteJSON(&FeedCmd{"subscribe", args})
}

func (f *PublicFeed) Unsubscribe(args interface{}) error {
	return f.WriteJSON(&FeedCmd{"unsubscribe", args})
}

func listenOn(reader io.Reader, feedChan chan *FeedCmd) {
	dec := json.NewDecoder(reader)

	for {
		feedData := &FeedCmd{}

		if err := dec.Decode(feedData); err != nil {
			log.Println(err)
			close(feedChan)
			break
		}

		feedChan <- feedData
	}
}
