package feed

import (
	"crypto/tls"
	"encoding/json"
)

type FeedCmd struct {
	Cmd  string      `json:"cmd"`
	Args interface{} `json:"args"`
}

type Feed struct {
	Service, SessionKey string
	Conn                *tls.Conn
}

func newFeed(address, service, sessionKey string) (*Feed, error) {
	if conn, err := tls.Dial("tcp", address, nil); err != nil {
		return nil, err
	} else {
		f := &Feed{Service: service, SessionKey: sessionKey, Conn: conn}
		return f, nil
	}
}

func (f *Feed) WriteJSON(jsonDoc interface{}) error {
	if jsonData, err := json.Marshal(jsonDoc); err != nil {
		return err
	} else {
		_, err = f.Conn.Write(append(jsonData, '\n'))
		return err
	}
}

func (f *Feed) Login() error {
	authData := &map[string]string{"session_key": f.SessionKey, "service": f.Service}
	loginCmd := &FeedCmd{"login", authData}
	return f.WriteJSON(loginCmd)
}

func (f *Feed) Close() error {
	return f.Conn.Close()
}
