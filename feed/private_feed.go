package feed

type PrivateFeed struct {
	*Feed
}

func NewPrivateFeed(address, service, sessionKey string) (*PrivateFeed, error) {
	f, err := newFeed(address, service, sessionKey)
	if err != nil {
		return nil, err
	}

	privf := &PrivateFeed{Feed: f}

	if err = privf.Login(); err != nil {
		return nil, err
	}

	return privf, nil
}
