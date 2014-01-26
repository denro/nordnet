package feed

type PublicFeed struct {
	*Feed
}

func NewPublicFeed(address, service, sessionKey string) (*PublicFeed, error) {
	f, err := NewFeed(address, service, sessionKey)
	if err != nil {
		return nil, err
	}

	pubf := &PublicFeed{Feed: f}
	if err = pubf.Login(); err != nil {
		return nil, err
	}

	return pubf, nil
}

func (f *PublicFeed) Subscribe(args interface{}) error {
	return f.Write(&FeedMsg{Cmd: "subscribe", Args: args})
}

func (f *PublicFeed) Unsubscribe(args interface{}) error {
	return f.Write(&FeedMsg{Cmd: "unsubscribe", Args: args})
}
