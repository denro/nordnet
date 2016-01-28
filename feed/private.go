package feed

import (
	"encoding/json"

	"github.com/denro/nordnet/util/models"
)

type PrivateFeed struct {
	*Feed
}

func NewPrivateFeed(address string) (*PrivateFeed, error) {
	f, err := newFeed(address)
	if err != nil {
		return nil, err
	}
	return &PrivateFeed{f}, nil
}

// Order data section in the private message
type PrivateOrder models.Order

// Trade data section in the private message
type PrivateTrade models.Trade

// Represents the messages sent on the private feed
type PrivateMsg FeedMsg

// Implements the Unmarshaler interface
// decodes the json into proper data types depending on the type field
func (pm *PrivateMsg) UnmarshalJSON(b []byte) (err error) {
	rawMsg := rawMsg{} // to avoid endless recursion below
	if err = json.Unmarshal(b, &rawMsg); err != nil {
		return
	}

	*pm = PrivateMsg{Type: rawMsg.Type}

	switch rawMsg.Type {
	case heartbeatType:
		pm.Data = struct{}{}
	case orderType:
		order := PrivateOrder{}
		if err = json.Unmarshal(rawMsg.Data, &order); err != nil {
			return
		}
		pm.Data = order
	case tradeType:
		trade := PrivateTrade{}
		if err = json.Unmarshal(rawMsg.Data, &trade); err != nil {
			return
		}
		pm.Data = trade
	}

	return
}

// Starts reading from the connection, returns channels for reading the messages and errors
func (pf *PrivateFeed) Dispatch(msgChan chan<- *PrivateMsg, errChan chan<- error) {
	go func(d *json.Decoder, mc chan<- *PrivateMsg, ec chan<- error) {
		var (
			pMsg *PrivateMsg
			err  error
		)

		for {
			pMsg = new(PrivateMsg)
			if err = d.Decode(pMsg); err != nil {
				ec <- err
			}
			msgChan <- pMsg
		}
	}(pf.decoder, msgChan, errChan)

	return
}
