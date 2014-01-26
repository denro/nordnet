package feed

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
)

type FeedMsg struct {
	Cmd  string      `json:"cmd,omitempty"`
	Type string      `json:"type,omitempty"`
	Args interface{} `json:"args,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type FeedArgs struct {
	T string `json:"t"`
	M int64  `json:"m"`
	I string `json:"i"`
}

type Feed struct {
	Service, SessionKey string
	Conn                *tls.Conn
}

func NewFeed(address, service, sessionKey string) (*Feed, error) {
	conn, err := tls.Dial("tcp", address, nil)
	if err != nil {
		return nil, err
	}

	return &Feed{service, sessionKey, conn}, nil
}

func (f *Feed) Write(any interface{}) error {
	json, err := json.Marshal(any)
	if err != nil {
		return err
	}

	_, err = f.Conn.Write(append(json, '\n'))
	return err
}

func (f *Feed) Login() error {
	authArgs := &map[string]string{"session_key": f.SessionKey, "service": f.Service}
	loginCmd := &FeedMsg{Cmd: "login", Args: authArgs}
	return f.Write(loginCmd)
}

func (f *Feed) Close() error {
	return f.Conn.Close()
}

func (f *Feed) DispatchListener(feedChan chan<- *FeedMsg) {
	go listenOn(f.Conn, feedChan)
}

func listenOn(reader io.Reader, feedChan chan<- *FeedMsg) {
	dec := json.NewDecoder(reader)

	for {
		feedData := &FeedMsg{}

		if err := dec.Decode(feedData); err != nil {
			log.Println(err)
			close(feedChan)
			break
		}

		feedChan <- feedData
	}
}
