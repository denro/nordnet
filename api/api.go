/*
	Package api includes the HTTP client used to access the REST JSON API.

	Information about specific endpoints and their parameters can be found at: https://api.test.nordnet.se/api-docs/index.html
*/
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

const (
	NNBASEURL    = `https://api.nordnet.se/next`
	NNSERVICE    = `NEXTAPI`
	NNAPIVERSION = `2`
)

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
func (c *APIClient) SystemStatus() (*SystemStatusResp, error) {
	res := &SystemStatusResp{}

	if err := c.Perform("GET", "", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list of accounts that the user has access to.
func (c *APIClient) Accounts() (*AccountsResp, error) {
	res := &AccountsResp{}

	if err := c.Perform("GET", "accounts", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// The account summary gives details of the account.
func (c *APIClient) Account(accountno int64) (*AccountResp, error) {
	res := &AccountResp{}

	path := fmt.Sprintf("accounts/%d", accountno)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Information about the currency ledgers of an account.
func (c *APIClient) AccountLedgers(accountno int64) (*AccountLedgersResp, error) {
	res := &AccountLedgersResp{}

	path := fmt.Sprintf("accounts/%d/ledgers", accountno)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all orders beloning to an account.
func (c *APIClient) AccountOrders(accountno int64, params *Params) (*AccountOrdersResp, error) {
	res := &AccountOrdersResp{}

	path := fmt.Sprintf("accounts/%d/orders", accountno)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Enter a new order, market_id + identifier is the identifier of the tradable.
func (c *APIClient) CreateOrder(accountno int64, params *Params) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("accounts/%d/orders", accountno)
	if err := c.Perform("POST", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Activate an inactive order. Please note that it is not possible to deactivate an order. The order must be entered as inactive.
func (c *APIClient) ActivateOrder(accountno int64, orderId int64) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("accounts/%d/orders/%d/activate", accountno, orderId)
	if err := c.Perform("PUT", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Modify price and or volume on an order.
func (c *APIClient) UpdateOrder(accountno int64, orderId int64, params *Params) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("accounts/%d/orders/%d", accountno, orderId)
	if err := c.Perform("PUT", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Delete an order.
func (c *APIClient) DeleteOrder(accountno int64, orderId int64) (*OrderResp, error) {
	res := &OrderResp{}

	path := fmt.Sprintf("accounts/%d/orders/%d", accountno, orderId)
	if err := c.Perform("DELETE", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list of all positions of the account.
func (c *APIClient) AccountPositions(accountno int64) (*AccountPositionsResp, error) {
	res := &AccountPositionsResp{}

	path := fmt.Sprintf("accounts/%d/positions", accountno)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all trades belonging to an account.
func (c *APIClient) AccountTrades(accountno int64, params *Params) (*AccountTradesResp, error) {
	res := &AccountTradesResp{}

	path := fmt.Sprintf("accounts/%d/trades", accountno)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get a list of all countries in the system. Please note that trading is not available everywhere.
func (c *APIClient) Countries() (*CountriesResp, error) {
	res := &CountriesResp{}

	if err := c.Perform("GET", "countries", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup one or more countries by country code. Multiple countries can be queried at the same time by comma separating the country codes.
// TODO: Merge with Countries call above
func (c *APIClient) LookupCountries(countries string) (*CountriesResp, error) {
	res := &CountriesResp{}

	path := fmt.Sprintf("countries/%s", countries)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list indicators that the user has access to.
func (c *APIClient) Indicators() (*IndicatorsResp, error) {
	res := &IndicatorsResp{}

	if err := c.Perform("GET", "indicators", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns info of one or more indicators.
// TODO: Merge with Indicators call above
func (c *APIClient) LookupIndicators(indicators string) (*IndicatorsResp, error) {
	res := &IndicatorsResp{}

	path := fmt.Sprintf("indicators/%s", indicators)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Free text search. A list of instruments is returned.
func (c *APIClient) SearchInstruments(params *Params) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	if err := c.Perform("GET", "instruments", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get one or more instruments, the instrument id is used as key
func (c *APIClient) Instruments(ids string) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	path := fmt.Sprintf("instruments/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list of leverage instruments that have the current instrument as underlying. Leverage instruments is for example warrants and ETF:s. To get all valid filters for the current underlying please use "Get leverages filters". The filters can be used to narrow the search. If "Get leverages filters" is used to fill comboboxes the same filters can be applied on the that call to hide filter cominations that are not valid. Multiple filters can be applied.
func (c *APIClient) InstrumentLeverages(id int64, params *Params) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	path := fmt.Sprintf("instruments/%d/leverages", id)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns valid filter values. Can be used to fill comboboxes in clients to filter leverages results. The same filters can be applied on this request to exclude invalid filter combinations.
func (c *APIClient) InstrumentLeverageFilters(id int64, params *Params) (*InstrumentLeverageFilterResp, error) {
	res := &InstrumentLeverageFilterResp{}

	path := fmt.Sprintf("instruments/%d/leverages/filters", id)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list of call/put option pairs. They are balanced on strike price. In order to find underlyings with options use "Get underlyings". To get available expiration dates use "Get option pair filters".
func (c *APIClient) InstrumentOptionPairs(id int64, params *Params) (*InstrumentOptionPairsResp, error) {
	res := &InstrumentOptionPairsResp{}

	path := fmt.Sprintf("instruments/%d/option_pairs", id)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns valid filter values. Can be used to fill comboboxes in clients to filter options pair results. The same filters can be applied on this request to exclude invalid filter combinations.
func (c *APIClient) InstrumentOptionPairFilters(id int64, params *Params) (*InstrumentOptionPairFiltersResp, error) {
	res := &InstrumentOptionPairFiltersResp{}

	path := fmt.Sprintf("instruments/%d/option_pairs/filters", id)
	if err := c.Perform("GET", path, params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup specfic instrument with prededfined fields. Please note that this is not a search, only exact matches is returned.
func (c *APIClient) InstrumentLookup(lookupType string, lookup string) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	path := fmt.Sprintf("instruments/lookup/%s/%s", lookupType, lookup)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all instrument sectors or the ones matching the group crtieria
func (c *APIClient) InstrumentSectors(params *Params) (*InstrumentSectorsResp, error) {
	res := &InstrumentSectorsResp{}

	if err := c.Perform("GET", "instruments/sectors", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get one or more sectors
func (c *APIClient) InstrumentSector(sectors string) (*InstrumentSectorsResp, error) {
	res := &InstrumentSectorsResp{}

	path := fmt.Sprintf("instruments/sectors/%s", sectors)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all instrument types. Please note that these types is used for both instrument_type and instrument_group_type.
func (c *APIClient) InstrumentTypes() (*InstrumentTypesResp, error) {
	res := &InstrumentTypesResp{}

	if err := c.Perform("GET", "instruments/types", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get info of one orde more instrument type.
func (c *APIClient) InstrumentType(instrumentType string) (*InstrumentTypesResp, error) {
	res := &InstrumentTypesResp{}

	path := fmt.Sprintf("instruments/types/%s", instrumentType)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get instruments that are underlyings for a specific type of instruments. The query can return instrument that have option derivatives or leverage derivatives. Warrants are included in the leverage derivatives.
func (c *APIClient) InstrumentUnderlyings(derivateType string, currency string) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	path := fmt.Sprintf("instruments/underlyings/%s/%s", derivateType, currency)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all instrument lists
func (c *APIClient) Lists() (*ListsResp, error) {
	res := &ListsResp{}

	if err := c.Perform("GET", "lists", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all instruments in a list.
func (c *APIClient) List(id int64) (*InstrumentsResp, error) {
	res := &InstrumentsResp{}

	path := fmt.Sprintf("lists/%d", id)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Before any other of the services (except for the system info request) can be called the user must login. The username, password and phrase must be sent encrypted.
func (c *APIClient) Login() (*LoginResp, error) {
	res := &LoginResp{}

	c.RLock()
	params := &Params{"auth": c.Credentials, "service": c.Service}
	c.RUnlock()

	if err := c.Perform("POST", "login", params, res); err != nil {
		return nil, err
	}

	c.Lock()
	c.SessionKey = res.SessionKey
	c.Unlock()

	return res, nil
}

// Invalidates the session.
func (c *APIClient) Logout() (*LogoutResp, error) {
	res := &LogoutResp{}

	if err := c.Perform("DELETE", "login", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// If the application needs to keep the session alive the session can be touched. Note the basic auth header field must be set as for all other calls. All calls to any REST service is touching the session. So touching the session manually is only needed if no other calls are done during the session timeout interval.

func (c *APIClient) Touch() (*TouchResp, error) {
	res := &TouchResp{}

	if err := c.Perform("PUT", "login", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

//Get all tradable markets. Market 80 is the smart order market. Instruments that can be traded on 2 or more markets gets a tradable on the smart order market. Orders entered with the smart order tradable get smart order routed with the current Nordnet best execution policy.
func (c *APIClient) Markets() (*MarketsResp, error) {
	res := &MarketsResp{}

	if err := c.Perform("GET", "markets", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Lookup one or more markets by market_id. Multiple market can be queried at the same time by comma separating the market_ids. Market 80 is the smart order market. Instruments that can be traded on 2 or more markets gets a tradable on the smart order market. Orders entered with the smart order tradable get smart order routed with the current Nordnet best execution policy.
func (c *APIClient) Market(ids string) (*MarketsResp, error) {
	res := &MarketsResp{}

	path := fmt.Sprintf("markets/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Search for news. If no search field is used the last news available to the user is returned.
func (c *APIClient) SearchNews(params *Params) (*NewsResp, error) {
	res := &NewsResp{}

	if err := c.Perform("GET", "news", params, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Show one or more news items.
// Search for news. If no search field is used the last news available to the user is returned.
func (c *APIClient) News(ids string) (*NewsItemsResp, error) {
	res := &NewsItemsResp{}

	path := fmt.Sprintf("news/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Returns a list of news sources the user has access to
func (c *APIClient) NewsSources() (*NewsSourcesResp, error) {
	res := &NewsSourcesResp{}

	if err := c.Perform("GET", "news_sources", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get realtime data access. This applies to the access on the feeds. If the market is missing the user don't have realtime access on that market.
func (c *APIClient) RealtimeAccess() (*RealtimeAccessResp, error) {
	res := &RealtimeAccessResp{}

	if err := c.Perform("GET", "realtime_access", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all ticksize tables.
func (c *APIClient) TickSizes() (*TickSizesResp, error) {
	res := &TickSizesResp{}

	if err := c.Perform("GET", "tick_sizes", nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get one or more ticksize tables.
func (c *APIClient) TickSize(ids string) (*TickSizesResp, error) {
	res := &TickSizesResp{}

	path := fmt.Sprintf("tick_sizes/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get trading calender and allowed trading types for one or more tradable.
func (c *APIClient) TradableInfo(ids string) (*TradableInfoResp, error) {
	res := &TradableInfoResp{}

	path := fmt.Sprintf("tradables/info/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Can be used for populating instrument price graphs for today. Resolution is one minute.
func (c *APIClient) TradableIntraday(ids string) (*TradableIntradayResp, error) {
	res := &TradableIntradayResp{}

	path := fmt.Sprintf("tradables/intraday/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Get all public trades (all trades done on the marketplace) beloning to one ore more tradable.
func (c *APIClient) TradableTrades(ids string) (*TradableTradesResp, error) {
	res := &TradableTradesResp{}

	path := fmt.Sprintf("tradables/trades/%s", ids)
	if err := c.Perform("GET", path, nil, res); err != nil {
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
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.RLock()
	if c.SessionKey != "" {
		req.SetBasicAuth(c.SessionKey, c.SessionKey)
	}
	c.RUnlock()

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
