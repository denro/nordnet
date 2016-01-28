/*
	Package api includes the HTTP client used to access the REST JSON API.

	Information about specific endpoints and their parameters can be found at: https://api.test.nordnet.se/api-docs/index.html
*/
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	. "github.com/denro/nordnet/util/models"
)

const (
	NNBASEURL    = `https://api.test.nordnet.se/next`
	NNSERVICE    = `NEXTAPI`
	NNAPIVERSION = `2`
)

var (
	TooManyRequestsError = errors.New("Too Many Requests, please wait for 10 seconds before trying again")
)

// Error type for errors returned by the API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIError implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

// Represents the options available for various methods.
type Params map[string]string

// APIClient provides all API-endpoints available as methods.
type APIClient struct {
	URL, Service, Version, Credentials, SessionKey string
	ExpiresAt, LastUsageAt                         time.Time

	http.Client
	sync.RWMutex
}

// Constructor function takes the credentials string produced by the util package.
func NewAPIClient(credentials string) *APIClient {
	return &APIClient{
		URL:         NNBASEURL,
		Service:     NNSERVICE,
		Version:     NNAPIVERSION,
		Credentials: credentials,
	}
}

// Information about the system status can be retrieved by this HTTP request. This is the only service that can be called without authentication.
func (c *APIClient) SystemStatus() (res *SystemStatus, err error) {
	res = &SystemStatus{}
	err = c.Perform("GET", "", nil, res)
	return
}

// Returns a list of accounts that the user has access to.
func (c *APIClient) Accounts() (res []Account, err error) {
	res = []Account{}
	err = c.Perform("GET", "accounts", nil, &res)
	return
}

// The account summary gives details of the account.
func (c *APIClient) Account(accountno int64) (res *AccountInfo, err error) {
	res = &AccountInfo{}
	err = c.Perform("GET", fmt.Sprintf("accounts/%d", accountno), nil, res)
	return
}

// Information about the currency ledgers of an account.
func (c *APIClient) AccountLedgers(accountno int64) (res []LedgerInformation, err error) {
	res = []LedgerInformation{}
	err = c.Perform("GET", fmt.Sprintf("accounts/%d/ledgers", accountno), nil, &res)
	return
}

// Get all orders beloning to an account.
func (c *APIClient) AccountOrders(accountno int64, params *Params) (res []Order, err error) {
	res = []Order{}
	err = c.Perform("GET", fmt.Sprintf("accounts/%d/orders", accountno), params, &res)
	return
}

// Enter a new order, market_id + identifier is the identifier of the tradable.
func (c *APIClient) CreateOrder(accountno int64, params *Params) (res *OrderReply, err error) {
	res = &OrderReply{}
	err = c.Perform("POST", fmt.Sprintf("accounts/%d/orders", accountno), params, res)
	return
}

// Activate an inactive order. Please note that it is not possible to deactivate an order. The order must be entered as inactive.
func (c *APIClient) ActivateOrder(accountno int64, orderId int64) (res *OrderReply, err error) {
	res = &OrderReply{}
	err = c.Perform("PUT", fmt.Sprintf("accounts/%d/orders/%d/activate", accountno, orderId), nil, res)
	return
}

// Modify price and or volume on an order.
func (c *APIClient) UpdateOrder(accountno int64, orderId int64, params *Params) (res *OrderReply, err error) {
	res = &OrderReply{}
	err = c.Perform("PUT", fmt.Sprintf("accounts/%d/orders/%d", accountno, orderId), params, res)
	return
}

// Delete an order.
func (c *APIClient) DeleteOrder(accountno int64, orderId int64) (res *OrderReply, err error) {
	res = &OrderReply{}
	err = c.Perform("DELETE", fmt.Sprintf("accounts/%d/orders/%d", accountno, orderId), nil, res)
	return
}

// Returns a list of all positions of the account.
func (c *APIClient) AccountPositions(accountno int64) (res []Position, err error) {
	res = []Position{}
	err = c.Perform("GET", fmt.Sprintf("accounts/%d/positions", accountno), nil, &res)
	return
}

// Get all trades belonging to an account.
func (c *APIClient) AccountTrades(accountno int64, params *Params) (res []Trade, err error) {
	res = []Trade{}
	err = c.Perform("GET", fmt.Sprintf("accounts/%d/trades", accountno), params, &res)
	return
}

// Get a list of all countries in the system. Please note that trading is not available everywhere.
func (c *APIClient) Countries() (res []Country, err error) {
	res = []Country{}
	c.Perform("GET", "countries", nil, &res)
	return
}

// Lookup one or more countries by country code. Multiple countries can be queried at the same time by comma separating the country codes.
// TODO: Merge with Countries call above?
func (c *APIClient) LookupCountries(countries string) (res []Country, err error) {
	res = []Country{}
	err = c.Perform("GET", fmt.Sprintf("countries/%s", countries), nil, &res)
	return
}

// Returns a list indicators that the user has access to.
func (c *APIClient) Indicators() (res []Indicator, err error) {
	res = []Indicator{}
	err = c.Perform("GET", "indicators", nil, &res)
	return
}

// Returns info of one or more indicators.
// TODO: Merge with Indicators call above?
func (c *APIClient) LookupIndicators(indicators string) (res []Indicator, err error) {
	res = []Indicator{}
	err = c.Perform("GET", fmt.Sprintf("indicators/%s", indicators), nil, &res)
	return
}

// Free text search. A list of instruments is returned.
func (c *APIClient) SearchInstruments(params *Params) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", "instruments", params, &res)
	return
}

// Get one or more instruments, the instrument id is used as key
func (c *APIClient) Instruments(ids string) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", fmt.Sprintf("instruments/%s", ids), nil, &res)
	return
}

// Returns a list of leverage instruments that have the current instrument as underlying. Leverage instruments is for example warrants and ETF:s. To get all valid filters for the current underlying please use "Get leverages filters". The filters can be used to narrow the search. If "Get leverages filters" is used to fill comboboxes the same filters can be applied on the that call to hide filter cominations that are not valid. Multiple filters can be applied.
func (c *APIClient) InstrumentLeverages(id int64, params *Params) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", fmt.Sprintf("instruments/%d/leverages", id), params, &res)
	return
}

// Returns valid filter values. Can be used to fill comboboxes in clients to filter leverages results. The same filters can be applied on this request to exclude invalid filter combinations.
func (c *APIClient) InstrumentLeverageFilters(id int64, params *Params) (res *LeverageFilter, err error) {
	res = &LeverageFilter{}
	err = c.Perform("GET", fmt.Sprintf("instruments/%d/leverages/filters", id), params, res)
	return
}

// Returns a list of call/put option pairs. They are balanced on strike price. In order to find underlyings with options use "Get underlyings". To get available expiration dates use "Get option pair filters".
func (c *APIClient) InstrumentOptionPairs(id int64, params *Params) (res []OptionPair, err error) {
	res = []OptionPair{}
	err = c.Perform("GET", fmt.Sprintf("instruments/%d/option_pairs", id), params, &res)
	return
}

// Returns valid filter values. Can be used to fill comboboxes in clients to filter options pair results. The same filters can be applied on this request to exclude invalid filter combinations.
func (c *APIClient) InstrumentOptionPairFilters(id int64, params *Params) (res *OptionPairFilter, err error) {
	res = &OptionPairFilter{}
	err = c.Perform("GET", fmt.Sprintf("instruments/%d/option_pairs/filters", id), params, res)
	return
}

// Lookup specfic instrument with prededfined fields. Please note that this is not a search, only exact matches is returned.
func (c *APIClient) InstrumentLookup(lookupType string, lookup string) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", fmt.Sprintf("instruments/lookup/%s/%s", lookupType, lookup), nil, &res)
	return
}

// Get all instrument sectors or the ones matching the group crtieria
func (c *APIClient) InstrumentSectors(params *Params) (res []Sector, err error) {
	res = []Sector{}
	err = c.Perform("GET", "instruments/sectors", params, &res)
	return
}

// Get one or more sectors
func (c *APIClient) InstrumentSector(sectors string) (res []Sector, err error) {
	res = []Sector{}
	err = c.Perform("GET", fmt.Sprintf("instruments/sectors/%s", sectors), nil, &res)
	return
}

// Get all instrument types. Please note that these types is used for both instrument_type and instrument_group_type.
func (c *APIClient) InstrumentTypes() (res []InstrumentType, err error) {
	res = []InstrumentType{}
	err = c.Perform("GET", "instruments/types", nil, &res)
	return
}

// Get info of one orde more instrument type.
func (c *APIClient) InstrumentType(instrumentType string) (res []InstrumentType, err error) {
	res = []InstrumentType{}
	err = c.Perform("GET", fmt.Sprintf("instruments/types/%s", instrumentType), nil, &res)
	return
}

// Get instruments that are underlyings for a specific type of instruments. The query can return instrument that have option derivatives or leverage derivatives. Warrants are included in the leverage derivatives.
func (c *APIClient) InstrumentUnderlyings(derivateType string, currency string) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", fmt.Sprintf("instruments/underlyings/%s/%s", derivateType, currency), nil, &res)
	return
}

// Get all instrument lists
func (c *APIClient) Lists() (res []List, err error) {
	res = []List{}
	err = c.Perform("GET", "lists", nil, &res)
	return
}

// Get all instruments in a list.
func (c *APIClient) List(id int64) (res []Instrument, err error) {
	res = []Instrument{}
	err = c.Perform("GET", fmt.Sprintf("lists/%d", id), nil, &res)
	return
}

// Before any other of the services (except for the system info request) can be called the user must login. The username, password and phrase must be sent encrypted.
// TODO: move the params into function arguments since its only used here?
func (c *APIClient) Login() (res *Login, err error) {
	res = &Login{}

	c.RLock()
	params := &Params{"auth": c.Credentials, "service": c.Service}
	c.RUnlock()

	err = c.Perform("POST", "login", params, res)

	c.Lock()
	c.SessionKey = res.SessionKey
	c.Unlock()

	return
}

// Invalidates the session.
func (c *APIClient) Logout() (res *LoggedInStatus, err error) {
	res = &LoggedInStatus{}
	err = c.Perform("DELETE", "login", nil, res)
	return
}

// If the application needs to keep the session alive the session can be touched. Note the basic auth header field must be set as for all other calls. All calls to any REST service is touching the session. So touching the session manually is only needed if no other calls are done during the session timeout interval.
func (c *APIClient) Touch() (res *LoggedInStatus, err error) {
	res = &LoggedInStatus{}
	err = c.Perform("PUT", "login", nil, res)
	return
}

//Get all tradable markets. Market 80 is the smart order market. Instruments that can be traded on 2 or more markets gets a tradable on the smart order market. Orders entered with the smart order tradable get smart order routed with the current Nordnet best execution policy.
func (c *APIClient) Markets() (res []Market, err error) {
	res = []Market{}
	err = c.Perform("GET", "markets", nil, &res)
	return
}

// Lookup one or more markets by market_id. Multiple market can be queried at the same time by comma separating the market_ids. Market 80 is the smart order market. Instruments that can be traded on 2 or more markets gets a tradable on the smart order market. Orders entered with the smart order tradable get smart order routed with the current Nordnet best execution policy.
func (c *APIClient) Market(ids string) (res []Market, err error) {
	res = []Market{}
	err = c.Perform("GET", fmt.Sprintf("markets/%s", ids), nil, &res)
	return
}

// Search for news. If no search field is used the last news available to the user is returned.
func (c *APIClient) SearchNews(params *Params) (res []NewsPreview, err error) {
	res = []NewsPreview{}
	err = c.Perform("GET", "news", params, &res)
	return
}

// Show one or more news items.
// Search for news. If no search field is used the last news available to the user is returned.
func (c *APIClient) News(ids string) (res []NewsItem, err error) {
	res = []NewsItem{}
	err = c.Perform("GET", fmt.Sprintf("news/%s", ids), nil, &res)
	return res, nil
}

// Returns a list of news sources the user has access to
func (c *APIClient) NewsSources() (res []NewsSource, err error) {
	res = []NewsSource{}
	err = c.Perform("GET", "news_sources", nil, &res)
	return
}

// Get realtime data access. This applies to the access on the feeds. If the market is missing the user don't have realtime access on that market.
func (c *APIClient) RealtimeAccess() (res []RealtimeAccess, err error) {
	res = []RealtimeAccess{}
	err = c.Perform("GET", "realtime_access", nil, &res)
	return
}

// Get all ticksize tables.
func (c *APIClient) TickSizes() (res []TicksizeTable, err error) {
	res = []TicksizeTable{}
	err = c.Perform("GET", "tick_sizes", nil, &res)
	return
}

// Get one or more ticksize tables.
func (c *APIClient) TickSize(ids string) (res []TicksizeTable, err error) {
	res = []TicksizeTable{}
	err = c.Perform("GET", fmt.Sprintf("tick_sizes/%s", ids), nil, &res)
	return
}

// Get trading calender and allowed trading types for one or more tradable.
func (c *APIClient) TradableInfo(ids string) (res []TradableInfo, err error) {
	res = []TradableInfo{}
	err = c.Perform("GET", fmt.Sprintf("tradables/info/%s", ids), nil, &res)
	return
}

// Can be used for populating instrument price graphs for today. Resolution is one minute.
func (c *APIClient) TradableIntraday(ids string) (res []IntradayGraph, err error) {
	res = []IntradayGraph{}
	err = c.Perform("GET", fmt.Sprintf("tradables/intraday/%s", ids), nil, &res)
	return
}

// Get all public trades (all trades done on the marketplace) beloning to one ore more tradable.
func (c *APIClient) TradableTrades(ids string) (res []PublicTrades, err error) {
	res = []PublicTrades{}
	err = c.Perform("GET", fmt.Sprintf("tradables/trades/%s", ids), nil, &res)
	return
}

func (c *APIClient) Perform(method, path string, params *Params, res interface{}) (err error) {
	reqURL, err := c.formatURL(path, params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, reqURL.String(), nil)
	if err != nil {
		return
	}

	resp, err := c.perform(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch resp.StatusCode {
	case 204:
		return
	case 400, 401, 404:
		errRes := APIError{}
		if err = json.Unmarshal(body, &errRes); err != nil {
			return
		}
		return errRes
	case 429:
		return TooManyRequestsError
	}

	if err = json.Unmarshal(body, res); err != nil {
		return
	}

	return
}

func (c *APIClient) perform(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.RLock()
	if c.SessionKey != "" {
		req.SetBasicAuth(c.SessionKey, c.SessionKey)
	}
	c.RUnlock()

	resp, err = c.Do(req)

	c.Lock()
	c.LastUsageAt = time.Now()
	c.Unlock()

	return
}

func (c *APIClient) formatURL(path string, params *Params) (*url.URL, error) {
	c.RLock()
	baseURL := fmt.Sprintf("%s/%s", c.URL, c.Version)
	c.RUnlock()

	if path != "" {
		baseURL += "/"
	}
	absURL := baseURL + path

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
