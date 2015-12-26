package feed

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var publicUnmarshalTests = []struct {
	json     string
	expected *PublicMsg
}{
	{
		`{"type":"heartbeat","data":{}}`,
		&PublicMsg{"heartbeat", struct{}{}},
	},
	{
		`{
			"type": "price",
			"data": {
				"i": "test",
				"m": 123,
				"trade_timestamp": 123,
				"tick_timestamp": 123,
				"bid": 1.1,
				"bid_volume": 1.1,
				"ask": 1.1,
				"ask_volume": 1.1,
				"close": 1.1,
				"high": 1.1,
				"last": 1.1,
				"last_volume": 1.1,
				"low": 1.1,
				"open": 1.1,
				"turnover": 1.1,
				"turnover_volume": 1.1,
				"ep": 1.1,
				"paired": 1.1,
				"imbalance": 1.1
			}
		}`,
		&PublicMsg{"price", PublicPrice{
			I:              "test",
			M:              123,
			TradeTimestamp: 123,
			TickTimestamp:  123,
			Bid:            1.1,
			BidVolume:      1.1,
			Ask:            1.1,
			AskVolume:      1.1,
			Close:          1.1,
			High:           1.1,
			Last:           1.1,
			LastVolume:     1.1,
			Low:            1.1,
			Open:           1.1,
			Turnover:       1.1,
			TurnoverVolume: 1.1,
			EP:             1.1,
			Paired:         1.1,
			Imbalance:      1.1,
		}},
	},
	{
		`{
			"type": "trade",
			"data": {
				"i": "test",
				"m": 123,
				"trade_timestamp": 123,
				"price": 1.1,
				"volume": 1.1,
				"broker_buying": "test",
				"broker_selling": "test",
				"trade_id": "test",
				"trade_type": "test"
			}
		}`,
		&PublicMsg{"trade", PublicTrade{
			I:              "test",
			M:              123,
			TradeTimestamp: 123,
			Price:          1.1,
			Volume:         1.1,
			BrokerBuying:   "test",
			BrokerSelling:  "test",
			TradeId:        "test",
			TradeType:      "test",
		}},
	},
	{
		`{
			"type": "depth",
			"data": {
				"i": "test",
				"m": 123,
				"tick_timestamp": 123,
				"bid1": 1.1,
				"bid_volume1": 1.1,
				"ask1": 1.1,
				"ask_volume1": 1.1,
				"bid2": 1.1,
				"bid_volume2": 1.1,
				"ask2": 1.1,
				"ask_volume2": 1.1,
				"bid3": 1.1,
				"bid_volume3": 1.1,
				"ask3": 1.1,
				"ask_volume3": 1.1,
				"bid4": 1.1,
				"bid_volume4": 1.1,
				"ask4": 1.1,
				"ask_volume4": 1.1,
				"bid5": 1.1,
				"bid_volume5": 1.1,
				"ask5": 1.1,
				"ask_volume5": 1.1
			}
		}`,
		&PublicMsg{"depth", PublicDepth{
			I:             "test",
			M:             123,
			TickTimestamp: 123,
			Bid1:          1.1,
			BidVolume1:    1.1,
			Ask1:          1.1,
			AskVolume1:    1.1,
			Bid2:          1.1,
			BidVolume2:    1.1,
			Ask2:          1.1,
			AskVolume2:    1.1,
			Bid3:          1.1,
			BidVolume3:    1.1,
			Ask3:          1.1,
			AskVolume3:    1.1,
			Bid4:          1.1,
			BidVolume4:    1.1,
			Ask4:          1.1,
			AskVolume4:    1.1,
			Bid5:          1.1,
			BidVolume5:    1.1,
			Ask5:          1.1,
			AskVolume5:    1.1,
		}},
	},
	{
		`{
			"type": "trading_status",
			"data": {
				"i": "test",
				"m": 123,
				"tick_timestamp": 123,
				"status": "test",
				"source_status": "test",
				"halted": "test"
			}
		}`,
		&PublicMsg{"trading_status", PublicTradingStatus{
			I:             "test",
			M:             123,
			TickTimestamp: 123,
			Status:        "test",
			SourceStatus:  "test",
			Halted:        "test",
		}},
	},
	{
		`{
			"type": "indicator",
			"data": {
				"i": "test",
				"m": "test",
				"tick_timestamp": 123,
				"high": 1.1,
				"low": 1.1,
				"last": 1.1,
				"close": 1.1
			}
		}`,
		&PublicMsg{"indicator", PublicIndicator{
			I:             "test",
			M:             "test",
			TickTimestamp: 123,
			High:          1.1,
			Low:           1.1,
			Last:          1.1,
			Close:         1.1,
		}},
	},
	{
		`{
			"type": "news",
			"data": {
				"itemid": "test",
				"lang": "test",
				"datetime": "test",
				"sourceid": "test",
				"headline": "test",
				"instruments": ["test"]
			}
		}`,
		&PublicMsg{"news", PublicNews{
			ItemId:      "test",
			Lang:        "test",
			Datetime:    "test",
			SourceId:    "test",
			Headline:    "test",
			Instruments: []string{"test"},
		}},
	},
}

func TestPublicMsgUnmarshalJSON(t *testing.T) {
	for _, tt := range publicUnmarshalTests {
		msg := &PublicMsg{}
		if err := json.Unmarshal([]byte(tt.json), msg); err != nil {
			t.Error(err)
		}
		assert.Equal(t, tt.expected, msg)
	}
}

var publicDispatchTests = []struct {
	json     string
	expected *PublicMsg
}{
	{
		`{"type":"heartbeat","data":{}}`,
		&PublicMsg{"heartbeat", struct{}{}},
	},
	{
		`{"type":"price","data":{}}`,
		&PublicMsg{"price", PublicPrice{}},
	},
	{
		`{"type":"trade","data":{}}`,
		&PublicMsg{"trade", PublicTrade{}},
	},
	{
		`{"type":"depth","data":{}}`,
		&PublicMsg{"depth", PublicDepth{}},
	},
	{
		`{"type":"indicator","data":{}}`,
		&PublicMsg{"indicator", PublicIndicator{}},
	},
	{
		`{"type":"news","data":{}}`,
		&PublicMsg{"news", PublicNews{}},
	},
	{
		`{"type":"trading_status","data":{}}`,
		&PublicMsg{"trading_status", PublicTradingStatus{}},
	},
}

func TestPublicFeedDispatch(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}
	feed := &PublicFeed{f}

	for _, tt := range publicDispatchTests {
		b.WriteString(tt.json + string('\n'))
	}

	msgChan, errChan := feed.Dispatch()

	for _, tt := range publicDispatchTests {
		select {
		case msg := <-msgChan:
			assert.Equal(t, tt.expected, msg)
		case err := <-errChan:
			t.Error(err)
		}
	}
}

var subscribeTests = []struct {
	args     interface{}
	expected string
}{
	{
		nil,
		`{"cmd":"subscribe","args":null}`,
	},
	{
		"test",
		`{"cmd":"subscribe","args":"test"}`,
	},
	{
		123,
		`{"cmd":"subscribe","args":123}`,
	},
	{
		map[string]interface{}{"some": "value"},
		`{"cmd":"subscribe","args":{"some":"value"}}`,
	},
	{
		PriceArgs{T: "price", I: "1869", M: 30},
		`{"cmd":"subscribe","args":{"t":"price","i":"1869","m":30}}`,
	},
	{
		DepthArgs{T: "depth", I: "1869", M: 30},
		`{"cmd":"subscribe","args":{"t":"depth","i":"1869","m":30}}`,
	},
	{
		TradeArgs{T: "trade", I: "1869", M: 30},
		`{"cmd":"subscribe","args":{"t":"trade","i":"1869","m":30}}`,
	},
	{
		TradingStatusArgs{T: "trading_status", I: "1869", M: 30},
		`{"cmd":"subscribe","args":{"t":"trading_status","i":"1869","m":30}}`,
	},
	{
		IndicatorArgs{T: "indicator", I: "SIX-IdX-DJI", M: "SIX"},
		`{"cmd":"subscribe","args":{"t":"indicator","i":"SIX-IdX-DJI","m":"SIX"}}`,
	},
	{
		NewsArgs{T: "news", S: 2, Delay: true},
		`{"cmd":"subscribe","args":{"t":"news","s":2,"delay":true}}`,
	},
	{
		NewsArgs{T: "news", S: 2},
		`{"cmd":"subscribe","args":{"t":"news","s":2}}`,
	},
}

func TestSubscribe(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}
	feed := &PublicFeed{f}

	for _, tt := range subscribeTests {
		b.Reset()
		feed.Subscribe(tt.args)
		assert.Equal(t, tt.expected+string('\n'), b.String())
	}
}

var unsubscribeTests = []struct {
	args     interface{}
	expected string
}{
	{
		nil,
		`{"cmd":"unsubscribe","args":null}`,
	},
	{
		"test",
		`{"cmd":"unsubscribe","args":"test"}`,
	},
	{
		123,
		`{"cmd":"unsubscribe","args":123}`,
	},
	{
		map[string]interface{}{"some": "value"},
		`{"cmd":"unsubscribe","args":{"some":"value"}}`,
	},
	{
		PriceArgs{T: "price", I: "1869", M: 30},
		`{"cmd":"unsubscribe","args":{"t":"price","i":"1869","m":30}}`,
	},
	{
		DepthArgs{T: "depth", I: "1869", M: 30},
		`{"cmd":"unsubscribe","args":{"t":"depth","i":"1869","m":30}}`,
	},
	{
		TradeArgs{T: "trade", I: "1869", M: 30},
		`{"cmd":"unsubscribe","args":{"t":"trade","i":"1869","m":30}}`,
	},
	{
		TradingStatusArgs{T: "trading_status", I: "1869", M: 30},
		`{"cmd":"unsubscribe","args":{"t":"trading_status","i":"1869","m":30}}`,
	},
	{
		IndicatorArgs{T: "indicator", I: "SIX-IdX-DJI", M: "SIX"},
		`{"cmd":"unsubscribe","args":{"t":"indicator","i":"SIX-IdX-DJI","m":"SIX"}}`,
	},
	{
		NewsArgs{T: "news", S: 2},
		`{"cmd":"unsubscribe","args":{"t":"news","s":2}}`,
	},
}

func TestUnsubscribe(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}
	feed := &PublicFeed{f}

	for _, tt := range unsubscribeTests {
		b.Reset()
		feed.Unsubscribe(tt.args)
		assert.Equal(t, tt.expected+string('\n'), b.String())
	}
}
