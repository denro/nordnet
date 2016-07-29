package feed

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
		&FeedCmd{"", nil},
		`{"cmd":"","args":null}`,
	},
	{
		&FeedCmd{"test", "test"},
		`{"cmd":"test","args":"test"}`,
	},
	{
		&FeedCmd{"test", 123},
		`{"cmd":"test","args":123}`,
	},
	{
		&FeedCmd{"test", map[string]interface{}{"some": "value"}},
		`{"cmd":"test","args":{"some":"value"}}`,
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
