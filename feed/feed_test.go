package feed

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeConnection struct {
	*bytes.Buffer
}

func (c *fakeConnection) Close() error {
	return nil
}

var writeTests = []struct {
	input    interface{}
	expected string
}{
	{
		&FeedCmd{"test", "test"},
		`{"cmd":"test","args":"test"}`,
	},
	{
		&FeedCmd{"login", LoginArgs{SessionKey: "ABC123"}},
		`{"cmd":"login","args":{"session_key":"ABC123"}}`,
	},
	{
		&FeedCmd{"login", LoginArgs{SessionKey: "ABC123", GetState: GetState{true, 2}}},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"deleted_orders":true,"days":2}}}`,
	},
	{
		&FeedCmd{"login", LoginArgs{SessionKey: "ABC123", GetState: GetState{true, 0}}},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"deleted_orders":true}}}`,
	},
	{
		&FeedCmd{"subscribe", PriceArgs{T: "price", I: "1869", M: 30}},
		`{"cmd":"subscribe","args":{"t":"price","i":"1869","m":30}}`,
	},
	{
		&FeedCmd{"subscribe", DepthArgs{T: "depth", I: "1869", M: 30}},
		`{"cmd":"subscribe","args":{"t":"depth","i":"1869","m":30}}`,
	},
	{
		&FeedCmd{"subscribe", TradeArgs{T: "trade", I: "1869", M: 30}},
		`{"cmd":"subscribe","args":{"t":"trade","i":"1869","m":30}}`,
	},
	{
		&FeedCmd{"subscribe", TradingStatusArgs{T: "trading_status", I: "1869", M: 30}},
		`{"cmd":"subscribe","args":{"t":"trading_status","i":"1869","m":30}}`,
	},
	{
		&FeedCmd{"subscribe", IndicatorArgs{T: "indicator", I: "SIX-IDX-DJI", M: "SIX"}},
		`{"cmd":"subscribe","args":{"t":"indicator","i":"SIX-IDX-DJI","m":"SIX"}}`,
	},
	{
		&FeedCmd{"subscribe", NewsArgs{T: "news", S: 2, Delay: true}},
		`{"cmd":"subscribe","args":{"t":"news","s":2,"delay":true}}`,
	},
	{
		&FeedCmd{"subscribe", NewsArgs{T: "news", S: 2}},
		`{"cmd":"subscribe","args":{"t":"news","s":2}}`,
	},
}

func TestWrite(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}

	for _, tt := range writeTests {
		b.Reset()
		if err := f.Write(tt.input); err != nil {
			t.Error(err)
		}
		assert.Equal(t, tt.expected+string('\n'), b.String())
	}
}

var loginTests = []struct {
	session  string
	getState interface{}
	expected string
}{
	{
		"ABC123",
		nil,
		`{"cmd":"login","args":{"session_key":"ABC123"}}`,
	},
	{
		"ABC123",
		&GetState{DeletedOrders: true},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"deleted_orders":true}}}`,
	},
	{
		"ABC123",
		&GetState{DeletedOrders: true, Days: 2},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"deleted_orders":true,"days":2}}}`,
	},
	{
		"ABC123",
		&map[string]interface{}{"deleted_orders": true},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"deleted_orders":true}}}`,
	},
	{
		"ABC123",
		&map[string]interface{}{"deleted_orders": true, "days": 2},
		`{"cmd":"login","args":{"session_key":"ABC123","get_state":{"days":2,"deleted_orders":true}}}`,
	},
}

func TestLogin(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}

	for _, tt := range loginTests {
		b.Reset()
		if err := f.Login(tt.session, tt.getState); err != nil {
			t.Error(err)
		}
		assert.Equal(t, tt.expected+string('\n'), b.String())
	}
}
