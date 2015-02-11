// Contains everything related to the public and private feeds
package feed

import (
	"crypto/tls"
	"encoding/json"
	"io"
)

// Used in the UnmarshalJSON implementations on PrivateFeed and PublicFeed
var (
	heartbeatType     = "heartbeat"
	orderType         = "order"
	tradeType         = "trade"
	priceType         = "price"
	depthType         = "depth"
	tradingStatusType = "trading_status"
	indicatorType     = "indicator"
	newsType          = "news"
)

// Used when sending feed commands
type FeedCmd struct {
	Cmd  string      `json:"cmd"`
	Args interface{} `json:"args"`
}

// Arguments for sending the login command
type LoginArgs struct {
	SessionKey string      `json:"session_key"`
	GetState   interface{} `json:"get_state,omitempty"`
}

// Arguments for getting orders and trades when logging in
type GetState struct {
	DeletedOrders bool  `json:"deleted_orders"`
	Days          int64 `json:"days,omitempty"`
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

// Used for receiving messages
type FeedMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Used in UnmarshalJSON overrides
type rawMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Represents the feed connection
type Feed struct {
	conn    io.ReadWriteCloser
	encoder *json.Encoder
	decoder *json.Decoder
}

// Returns a new Feed connected to the address specified
func NewFeed(address string, conf *tls.Config) (*Feed, error) {
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		return nil, err
	}

	return &Feed{conn, json.NewEncoder(conn), json.NewDecoder(conn)}, nil
}

// Feed implements the Writer interface
func (f *Feed) Write(any interface{}) error {
	return f.encoder.Encode(any)
}

// Feed implements the Closer interface
// closes the underlying conneciton
func (f *Feed) Close() error {
	return f.conn.Close()
}

// Send the login command with the specified session key
func (f *Feed) Login(session string, getState interface{}) error {
	return f.Write(&FeedCmd{Cmd: "login", Args: &LoginArgs{SessionKey: session, GetState: getState}})
}
