package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Params map[string]string

type APIClient struct {
	URL, Service, Credentials, SessionKey string
	ExpiresAt, LastUsageAt                time.Time

	http.Client
	sync.Mutex
}

func NewAPIClient(credentials string) *APIClient {
	return &APIClient{
		Credentials: credentials,
		Service:     "NEXTAPI",
		URL:         "https://api.nordnet.se/next",
	}
}

type SystemStatusResp struct {
	Timestamp     int64  `json:"timestamp"`
	ValidVersion  bool   `json:"valid_version"`
	SystemRunnnig bool   `json:"system_running"`
	SkipPhrase    bool   `json:"skip_phrase"`
	Message       string `json:"message"`
}

func (c *APIClient) SystemStatus() (*SystemStatusResp, error) {
	res := &SystemStatusResp{}

	if err := c.Perform("GET", "v1", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type feed struct {
	Port      int64  `json:"port"`
	Hostname  string `json:"hostname"`
	Encrypted bool   `json:"encrypted"`
}

type LoginResp struct {
	SessionKey  string `json:"session_key"`
	Environment string `json:"environment"`
	ExpiresIn   int64  `json:"expires_in"`
	PublicFeed  feed   `json:"public_feed"`
	PrivateFeed feed   `json:"private_feed"`
}

func (c *APIClient) Login() (*LoginResp, error) {
	res := &LoginResp{}

	c.Lock()
	params := &Params{"auth": c.Credentials, "service": c.Service}
	c.Unlock()

	if err := c.Perform("POST", "v1/login", params, res); err != nil {
		return nil, err
	}

	c.Lock()
	c.SessionKey = res.SessionKey
	c.Unlock()

	return res, nil
}

type LogoutResp struct {
	LoggedIn bool `json:"logged_in"`
}

func (c *APIClient) Logout() (*LogoutResp, error) {
	res := &LogoutResp{}

	c.Lock()
	path := fmt.Sprintf("v1/login/%s", c.SessionKey)
	c.Unlock()

	if err := c.Perform("DELETE", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type TouchResp struct {
	LoggedIn bool `json:"logged_in"`
}

func (c *APIClient) Touch() (*TouchResp, error) {
	res := &TouchResp{}

	c.Lock()
	path := fmt.Sprintf("v1/login/%s", c.SessionKey)
	c.Unlock()

	if err := c.Perform("PUT", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type RealtimeAccessResp []struct {
	MarketId string `json:"marketID"`
	Level    int64  `json:"level"`
}

func (c *APIClient) RealtimeAccess() (*RealtimeAccessResp, error) {
	res := &RealtimeAccessResp{}

	if err := c.Perform("GET", "v1/realtime_access", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type NewsSourcesResp []struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Level    string `json:"level"`
	SourceId int64  `json:"sourceid"`
	ImageURL string `json:"imageurl"`
}

func (c *APIClient) NewsSources() (*NewsSourcesResp, error) {
	res := &NewsSourcesResp{}

	if err := c.Perform("GET", "v1/news_sources", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type NewsItemsResp []struct {
	DateTime string `json:"datetime"`
	Headline string `json:"headline"`
	ItemId   int64  `json:"itemid"`
	SourceId int64  `json:"sourceid"`
	Type     string `json:"type"`
}

func (c *APIClient) NewsItems(params *Params) (*NewsItemsResp, error) {
	res := &NewsItemsResp{}

	if err := c.Perform("GET", "v1/news_items", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

type NewsItemResp struct {
	DateTime string `json:"datetime"`
	Body     string `json:"body"`
	Headline string `json:"headline"`
	ItemId   int64  `json:"itemid"`
	Lang     string `json:"lang"`
	Preamble string `json:"preamble"`
	SourceId int64  `json:"sourceid"`
	Type     string `json:"type"`
}

func (c *APIClient) NewsItem(newsItemId int64) (*NewsItemResp, error) {
	res := &NewsItemResp{}

	path := fmt.Sprintf("v1/news_items/%d", newsItemId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountsResp []struct {
	Id      string `json:"id"`
	Default string `json:"default"`
}

func (c *APIClient) Accounts() (*AccountsResp, error) {
	res := &AccountsResp{}

	if err := c.Perform("GET", "v1/accounts", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountResp struct {
	AccountSum      float64 `json:"accountSum,string"`
	FullMarketValue float64 `json:"fullMarketvalue,string"`
	TradingPower    float64 `json:"tradingPower,string"`
	AccountCurrency string  `json:"accountCurrency"`
}

func (c *APIClient) Account(accountId string) (*AccountResp, error) {
	res := &AccountResp{}

	path := fmt.Sprintf("v1/accounts/%s", accountId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountLedgersResp []struct {
	Currency      string  `json:"currency"`
	AccountSum    float64 `json:"accountSum,string"`
	AccountSumAcc float64 `json:"accountSumAcc,string"`
	AccIntCred    float64 `json:"accIntCred,string"`
	AccIntDeb     float64 `json:"accIntDeb,string"`
}

func (c *APIClient) AccountLedgers(accountId string) (*AccountLedgersResp, error) {
	res := &AccountLedgersResp{}

	path := fmt.Sprintf("v1/accounts/%s/ledgers", accountId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountPositionsResp []struct {
	AcqPrice       float64 `json:"acqPrice,string"`
	AcqPriceAcc    float64 `json:"acqPriceAcc,string"`
	PawnPercent    float64 `json:"pawnPercent,string"`
	Qty            float64 `json:"qty,string"`
	MarketValue    float64 `json:"marketValue,string"`
	MarketValueAcc float64 `json:"marketValueAcc,string"`

	Instrument struct {
		MainMarketId    int64   `json:"mainMarketId,string"`
		Identifier      string  `json:"identifier"`
		Type            string  `json:"type"`
		Currency        string  `json:"currency"`
		MainMarketPrice float64 `json:"mainMarketPrice,string"`
	} `json:"instrumentID"`
}

func (c *APIClient) AccountPositions(accountId string) (*AccountPositionsResp, error) {
	res := &AccountPositionsResp{}

	path := fmt.Sprintf("v1/accounts/%s/positions", accountId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountOrdersResp []struct {
	ExchangeOrderId string `json:"exchangeOrderID"`
	OrderId         int64  `json:"orderID"`

	ActivationCondition struct {
		Type string `json:"type"`
		Date string `json:"date"`

		Price struct {
			Value float64 `json:"value"`
			Curr  string  `json:"curr"`
		} `json:"price"`
	} `json:"activationCondition"`

	RegDate        int64  `json:"regdate"`
	PriceCondition string `json:"priceCondition"`

	Price struct {
		Value float64 `json:"value"`
		Curr  string  `json:"curr"`
	} `json:"price"`

	VolumeCondition string  `json:"volumeCondition"`
	Volume          float64 `json:"volume"`
	Side            string  `json:"side"`
	TradedVolume    float64 `json:"tradedVolume"`
	Accno           int64   `json:"accno"`

	InstrumentId struct {
		MarketId   int64  `json:"marketID"`
		Identifier string `json:"identifier"`
	} `json:"instrumentID"`

	Validity struct {
		Type string `json:"type"`
		Date string `json:"date"`
	} `json:"validity"`

	OrderState string `json:"orderState"`
	StatusText string `json:"statusText"`
}

func (c *APIClient) AccountOrders(accountId string) (*AccountOrdersResp, error) {
	res := &AccountOrdersResp{}

	path := fmt.Sprintf("v1/accounts/%s/orders", accountId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type AccountTradesResp []struct {
	SecurityTrade struct {
		Accno   string `json:"accno"`
		OrderId int64  `json:"orderID,string"`

		InstrumentId struct {
			MarketId   int64  `json:"marketID,string"`
			Identifier string `json:"identifier"`
		} `json:"instrumentID"`

		Volume    float64 `json:"volume,string"`
		TradeTime string  `json:"tradetime"`

		Price struct {
			Value float64 `json:"value,string"`
			Curr  string  `json:"curr"`
		} `json:"price"`

		Side    string `json:"side"`
		TradeId string `json:"tradeID"`
	} `json:"securityTrade"`
}

func (c *APIClient) AccountTrades(accountId string) (*AccountTradesResp, error) {
	res := &AccountTradesResp{}

	path := fmt.Sprintf("v1/accounts/%s/trades", accountId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type InstrumentsResp []struct {
	Type       string `json:"type"`
	LongName   string `json:"longname"`
	ShortName  string `json:"shortname"`
	MarketId   int64  `json:"marketID,string"`
	MarketName string `json:"marketname"`
	Country    string `json:"country"`
	IsInCode   string `json:"isinCode"`
	Identifier string `json:"identifier"`
	Currency   string `json:"currency"`
}

func (c *APIClient) Instruments(params *Params) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	if err := c.Perform("GET", "v1/instruments", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

type InstrumentResp struct {
	Type       string `json:"type"`
	LongName   string `json:"longname"`
	ShortName  string `json:"shortname"`
	MarketId   int64  `json:"marketID,string"`
	MarketName string `json:"marketname"`
	Country    string `json:"country"`
	IsInCode   string `json:"isinCode"`
	Identifier string `json:"identifier"`
	Currency   string `json:"currency"`
}

func (c *APIClient) Instrument(params *Params) (*InstrumentResp, error) {
	res := &InstrumentResp{}

	if err := c.Perform("GET", "v1/instruments", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ChartDataResp []struct {
	Timestamp string  `json:"timestamp"`
	Change    float64 `json:"change"`
	Volume    int64   `json:"volume"`
	Price     float64 `json:"price"`
}

func (c *APIClient) ChartData(params *Params) (*ChartDataResp, error) {
	res := &ChartDataResp{}

	if err := c.Perform("GET", "v1/chart_data", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ListsResp []struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

func (c *APIClient) Lists() (*ListsResp, error) {
	res := &ListsResp{}

	if err := c.Perform("GET", "v1/lists", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type ListResp []struct {
	ShortName  string `json:"shortname"`
	MarketId   int64  `json:"marketID,string"`
	Identifier string `json:"identifier"`
}

func (c *APIClient) List(listId int64) (*ListResp, error) {
	res := &ListResp{}

	path := fmt.Sprintf("v1/lists/%d", listId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MarketsReps []struct {
	Name       string `json:"name"`
	Country    string `json:"country"`
	MarketId   int64  `json:"marketID,string"`
	OrderTypes []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"ordertypes"`
}

func (c *APIClient) Markets() (*MarketsReps, error) {
	res := &MarketsReps{}

	if err := c.Perform("GET", "v1/markets", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type MarketTradingDaysResp []struct {
	Date        string `json:"date"`
	DisplayDate string `json:"display_date"`
}

func (c *APIClient) MarketTradingDays(marketId int64) (*MarketTradingDaysResp, error) {
	res := &MarketTradingDaysResp{}

	path := fmt.Sprintf("v1/markets/%d/trading_days", marketId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type IndicesResp []struct {
	Type     string `json:"type"`
	LongName string `json:"longname"`
	Source   string `json:"source"`
	Country  string `json:"country"`
	ImageURL string `json:"imageurl"`
	Id       string `json:"id"`
}

func (c *APIClient) Indices() (*IndicesResp, error) {
	res := &IndicesResp{}

	if err := c.Perform("GET", "v1/indices", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type TicksizesResp []struct {
	Tick     float64 `json:"tick"`
	Above    float64 `json:"above"`
	Decimals int64   `json:"decimals"`
}

func (c *APIClient) Ticksizes(ticksizeId int64) (*TicksizesResp, error) {
	res := &TicksizesResp{}

	path := fmt.Sprintf("v1/ticksizes/%d", ticksizeId)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type DerivateCountriesResp []string

func (c *APIClient) DerivateCountries(derType string) (*DerivateCountriesResp, error) {
	res := &DerivateCountriesResp{}

	path := fmt.Sprintf("v1/derivatives/%s", derType)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type DerivateUnderlyingsResp []struct {
	ShortName  string `json:"shortname"`
	MarketId   int64  `json:"marketID,string"`
	Identifier string `json:"identifier"`
}

func (c *APIClient) DerivateUnderlyings(derType, country string) (*DerivateUnderlyingsResp, error) {
	res := &DerivateUnderlyingsResp{}

	path := fmt.Sprintf("v1/derivatives/%s/underlyings/%s", derType, country)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type Derivatives []struct {
	ShortName   string  `json:"shortname"`
	Multipier   int64   `json:"multiplier,string"`
	StrikePrice float64 `json:"strikeprice,string"`
	MarketId    int64   `json:"marketID,string"`
	Identifier  string  `json:"identifier"`
	ExpiryDate  string  `json:"expirydate"`
	ExpiryType  string  `json:"expirytype"`
	Kind        string  `json:"kind"`
	Currency    string  `json:"currency"`
	CallPut     string  `json:"callPut"`
}

func (c *APIClient) Derivatives(derType string, params *Params) (*Derivatives, error) {
	res := &Derivatives{}

	path := fmt.Sprintf("v1/derivatives/%s/derivatives", derType)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

type RelatedMarketsResp []struct {
	MarketId   int64  `json:"marketID"`
	Identifier string `json:"identifier"`
}

func (c *APIClient) RelatedMarkets(params *Params) (*RelatedMarketsResp, error) {
	res := &RelatedMarketsResp{}

	if err := c.Perform("GET", "v1/related_markets", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

type OrderResp struct {
	OrderId     int64  `json:"orderID"`
	ResultCode  string `json:"resultCode"`
	OrderState  string `json:"orderState"`
	AccNo       int64  `json:"accNo"`
	ActionState string `json:"actionState"`
}

func (c *APIClient) CreateOrder(accountId string, params *Params) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("v1/accounts/%s/orders", accountId)
	if err := c.Perform("POST", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) UpdateOrder(accountId string, orderId int64, params *Params) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("v1/accounts/%s/orders/%d", accountId, orderId)
	if err := c.Perform("PUT", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) DeleteOrder(accountId string, orderId int64) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("v1/accounts/%s/orders/%d", accountId, orderId)
	if err := c.Perform("DELETE", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Perform(method, path string, params *Params, res interface{}) error {
	if reqURL, err := c.formatURL(path, params); err != nil {
		return err
	} else if req, err := http.NewRequest(method, reqURL.String(), nil); err != nil {
		return err
	} else if resp, err := c.perform(req); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		} else if err := json.Unmarshal(body, res); err != nil {
			return err
		}
	}

	return nil
}

func (c *APIClient) perform(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")

	c.Lock()
	if c.SessionKey != "" {
		req.SetBasicAuth(c.SessionKey, c.SessionKey)
	}
	c.Unlock()

	if resp, err := c.Do(req); err != nil {
		return nil, err
	} else {
		c.Lock()
		c.LastUsageAt = time.Now()
		c.Unlock()

		return resp, nil
	}
}

func (c *APIClient) formatURL(path string, params *Params) (*url.URL, error) {
	c.Lock()
	absURL := fmt.Sprintf("%s/%s", c.URL, path)
	c.Unlock()

	if reqURL, err := url.Parse(absURL); err != nil {
		return nil, err
	} else {
		if params != nil {
			reqQuery := reqURL.Query()
			for key, value := range *params {
				reqQuery.Set(key, value)
			}
			reqURL.RawQuery = reqQuery.Encode()
		}

		return reqURL, nil
	}
}
