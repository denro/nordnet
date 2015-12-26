package feed

import (
	"encoding/json"
)

type PublicFeed struct {
	*Feed
}

func NewPublicFeed(address string) (*PublicFeed, error) {
	f, err := newFeed(address)
	if err != nil {
		return nil, err
	}

	return &PublicFeed{f}, err
}

// Arguments for subscribing to price updates
type PriceArgs feedCmdArgs

// Arguments for subscribing to depth updates
type DepthArgs feedCmdArgs

// Arguments for subscribing to trade updates
type TradeArgs feedCmdArgs

// Arguments for subscribing to trading status updates
type TradingStatusArgs feedCmdArgs

type feedCmdArgs struct {
	T string `json:"t"`
	I string `json:"i"`
	M int64  `json:"m"`
}

// Arguments for subscribing to indicator updates
type IndicatorArgs struct {
	T string `json:"t"`
	I string `json:"i"`
	M string `json:"m"`
}

// Arguments for subscribing to news updates
type NewsArgs struct {
	T     string `json:"t"`
	S     int64  `json:"s"`
	Delay bool   `json:"delay,omitempty"`
}

// Sends the Subscribe command with the given args
func (f *PublicFeed) Subscribe(args interface{}) error {
	return f.Write(&FeedCmd{Cmd: "subscribe", Args: args})
}

// Sends the Unsubscribe command with the given args
func (f *PublicFeed) Unsubscribe(args interface{}) error {
	return f.Write(&FeedCmd{Cmd: "unsubscribe", Args: args})
}

// Price data section in the public message
type PublicPrice struct {
	I              string  `json:"i"`
	M              int64   `json:"m"`
	TradeTimestamp int64   `json:"trade_timestamp"`
	TickTimestamp  int64   `json:"tick_timestamp"`
	Bid            float64 `json:"bid"`
	BidVolume      float64 `json:"bid_volume"`
	Ask            float64 `json:"ask"`
	AskVolume      float64 `json:"ask_volume"`
	Close          float64 `json:"close"`
	High           float64 `json:"high"`
	Last           float64 `json:"last"`
	LastVolume     float64 `json:"last_volume"`
	Low            float64 `json:"low"`
	Open           float64 `json:"open"`
	Turnover       float64 `json:"turnover"`
	TurnoverVolume float64 `json:"turnover_volume"`
	EP             float64 `json:"ep"`
	Paired         float64 `json:"paired"`
	Imbalance      float64 `json:"imbalance"`
}

// Trade data section in the public message
type PublicTrade struct {
	I              string  `json:"i"`
	M              int64   `json:"m"`
	TradeTimestamp int64   `json:"trade_timestamp"`
	Price          float64 `json:"price"`
	Volume         float64 `json:"volume"`
	BrokerBuying   string  `json:"broker_buying"`
	BrokerSelling  string  `json:"broker_selling"`
	TradeId        string  `json:"trade_id"`
	TradeType      string  `json:"trade_type"`
}

// Depth data section in the public message
type PublicDepth struct {
	I             string  `json:"i"`
	M             int64   `json:"m"`
	TickTimestamp int64   `json:"tick_timestamp"`
	Bid1          float64 `json:"bid1"`
	BidVolume1    float64 `json:"bid_volume1"`
	Ask1          float64 `json:"ask1"`
	AskVolume1    float64 `json:"ask_volume1"`
	Bid2          float64 `json:"bid2"`
	BidVolume2    float64 `json:"bid_volume2"`
	Ask2          float64 `json:"ask2"`
	AskVolume2    float64 `json:"ask_volume2"`
	Bid3          float64 `json:"bid3"`
	BidVolume3    float64 `json:"bid_volume3"`
	Ask3          float64 `json:"ask3"`
	AskVolume3    float64 `json:"ask_volume3"`
	Bid4          float64 `json:"bid4"`
	BidVolume4    float64 `json:"bid_volume4"`
	Ask4          float64 `json:"ask4"`
	AskVolume4    float64 `json:"ask_volume4"`
	Bid5          float64 `json:"bid5"`
	BidVolume5    float64 `json:"bid_volume5"`
	Ask5          float64 `json:"ask5"`
	AskVolume5    float64 `json:"ask_volume5"`
}

// Trading Status data section in the public message
type PublicTradingStatus struct {
	I             string `json:"i"`
	M             int64  `json:"m"`
	TickTimestamp int64  `json:"tick_timestamp"`
	Status        string `json:"status"`
	SourceStatus  string `json:"source_status"`
	Halted        string `json:"halted"`
}

// Indicator data section in the public message
type PublicIndicator struct {
	I             string  `json:"i"`
	M             string  `json:"m"`
	TickTimestamp int64   `json:"tick_timestamp"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Last          float64 `json:"last"`
	Close         float64 `json:"close"`
}

// News data section in the public message
type PublicNews struct {
	ItemId      string   `json:"itemid"`
	Lang        string   `json:"lang"`
	Datetime    string   `json:"datetime"`
	SourceId    string   `json:"sourceid"`
	Headline    string   `json:"headline"`
	Instruments []string `json:"instruments"`
}

// Represents the messages sent on the public feed
type PublicMsg FeedMsg

// Implements the Unmarshaler interface
// decodes the json into proper data types depending on the type field
func (pm *PublicMsg) UnmarshalJSON(b []byte) (err error) {
	rawMsg := rawMsg{} // to avoid endless recursion below
	if err = json.Unmarshal(b, &rawMsg); err != nil {
		return
	}

	*pm = PublicMsg{Type: rawMsg.Type}

	switch rawMsg.Type {
	case heartbeatType:
		pm.Data = struct{}{}
	case priceType:
		price := PublicPrice{}
		if err = json.Unmarshal(rawMsg.Data, &price); err != nil {
			return
		}
		pm.Data = price
	case tradeType:
		trade := PublicTrade{}
		if err = json.Unmarshal(rawMsg.Data, &trade); err != nil {
			return
		}
		pm.Data = trade
	case depthType:
		depth := PublicDepth{}
		if err = json.Unmarshal(rawMsg.Data, &depth); err != nil {
			return
		}
		pm.Data = depth
	case tradingStatusType:
		tradingStatus := PublicTradingStatus{}
		if err = json.Unmarshal(rawMsg.Data, &tradingStatus); err != nil {
			return
		}
		pm.Data = tradingStatus
	case indicatorType:
		indicator := PublicIndicator{}
		if err = json.Unmarshal(rawMsg.Data, &indicator); err != nil {
			return
		}
		pm.Data = indicator
	case newsType:
		news := PublicNews{}
		if err = json.Unmarshal(rawMsg.Data, &news); err != nil {
			return
		}
		pm.Data = news
	}

	return
}

// Starts reading from the connection, returns channels for reading the messages and errors
func (pf *PublicFeed) Dispatch() (msgChan chan *PublicMsg, errChan chan error) {
	msgChan = make(chan *PublicMsg)
	errChan = make(chan error)

	go func(d *json.Decoder, mc chan<- *PublicMsg, ec chan<- error) {
		var (
			pMsg *PublicMsg
			err  error
		)

		for {
			pMsg = new(PublicMsg)
			if err = d.Decode(pMsg); err != nil {
				ec <- err
			}
			msgChan <- pMsg
		}
	}(pf.decoder, msgChan, errChan)

	return
}
