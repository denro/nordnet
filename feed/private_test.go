package feed

import (
	"bytes"
	"encoding/json"
	"github.com/denro/nordnet/util/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var privateUnmarshalTests = []struct {
	json     string
	expected *PrivateMsg
}{
	{
		`{
			"type":"heartbeat",
			"data":{}
		}`,
		&PrivateMsg{"heartbeat", struct{}{}},
	},
	{
		`{
			"type":"order",
			"data":{
				"accno":123,
				"order_id":123,
				"price":{
					"value":1.1,
					"currency":"test"
				},
				"volume":1.1,
				"tradable":{
					"identifier":"test",
					"market_id":123
				},
				"open_volume":1.1,
				"traded_volume":1.1,
				"side":"test",
				"modified":123,
				"reference":"test",
				"activation_condition":{
					"type":"test",
					"trailing_value":1.1,
					"trigger_value":1.1,
					"trigger_condition":"test"
				},
				"price_condition":"test",
				"volume_condition":"test",
				"validity":{
					"type":"test",
					"valid_until":123
				},
				"action_state":"test",
				"order_state":"test"
			}
		}`,
		&PrivateMsg{"order", PrivateOrder{
			Accno:               123,
			OrderId:             123,
			Price:               models.Amount{1.1, "test"},
			Volume:              1.1,
			Tradable:            models.TradableId{"test", 123},
			OpenVolume:          1.1,
			TradedVolume:        1.1,
			Side:                "test",
			Modified:            123,
			Reference:           "test",
			ActivationCondition: models.ActivationCondition{"test", 1.1, 1.1, "test"},
			PriceCondition:      "test",
			VolumeCondition:     "test",
			Validity:            models.Validity{"test", 123},
			ActionState:         "test",
			OrderState:          "test",
		}},
	},
	{
		`{
			"type": "trade",
			"data": {
				"accno": 123,
				"order_id": 123,
				"trade_id": "test",
				"tradable": {
					"identifier": "test",
					"market_id": 123
				},
				"price": {
					"value": 1.1,
					"currency": "test"
				},
				"volume": 1.1,
				"side": "test",
				"counterparty": "test",
				"tradetime": 123
			}
		}`,
		&PrivateMsg{"trade", PrivateTrade{
			Accno:        123,
			OrderId:      123,
			TradeId:      "test",
			Tradable:     models.TradableId{"test", 123},
			Price:        models.Amount{1.1, "test"},
			Volume:       1.1,
			Side:         "test",
			Counterparty: "test",
			Tradetime:    123,
		}},
	},
}

func TestPrivateMsgUnmarshalJSON(t *testing.T) {
	for _, tt := range privateUnmarshalTests {
		msg := &PrivateMsg{}
		if err := json.Unmarshal([]byte(tt.json), msg); err != nil {
			t.Error(err)
		}
		assert.Equal(t, tt.expected, msg)
	}
}

var privateDispatchTests = []struct {
	json     string
	expected *PrivateMsg
}{
	{
		`{"type":"heartbeat","data":{}}`,
		&PrivateMsg{"heartbeat", struct{}{}},
	},
	{
		`{"type":"trade","data":{}}`,
		&PrivateMsg{"trade", PrivateTrade{}},
	},
	{
		`{"type":"order","data":{}}`,
		&PrivateMsg{"order", PrivateOrder{}},
	},
}

func TestPrivateFeedDispatch(t *testing.T) {
	b := &fakeConnection{&bytes.Buffer{}}
	f := &Feed{b, json.NewEncoder(b), json.NewDecoder(b)}
	feed := &PrivateFeed{f}

	for _, tt := range privateDispatchTests {
		b.WriteString(tt.json + string('\n'))
	}

	msgChan, errChan := feed.Dispatch()

	for _, tt := range privateDispatchTests {
		select {
		case msg := <-msgChan:
			assert.Equal(t, tt.expected, msg)
		case err := <-errChan:
			t.Error(err)
		}
	}
}
