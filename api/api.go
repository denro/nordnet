package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Result map[string]interface{}
type ResultList []*Result

type APIClient struct {
	Service, Credentials, SessionKey string
	ExpiresAt, LastUsageAt           time.Time

	http.Client
}

func NewAPIClient(credentials, service string) *APIClient {
	return &APIClient{Credentials: credentials, Service: service}
}

func (c *APIClient) SystemStatus() (*Result, error) {
	res := &Result{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Login() (*Result, error) {
	res := &Result{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/login")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	reqQuery.Set("auth", c.Credentials)
	reqQuery.Set("service", c.Service)
	reqUrl.RawQuery = reqQuery.Encode()

	req, err := http.NewRequest("POST", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	err = c.DoAndUnmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (c *APIClient) Touch() (*Result, error) {
	res := &Result{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/login/%s", c.SessionKey)
	req, err := http.NewRequest("PUT", urlStr, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoAndUnmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Logout() (*Result, error) {
	res := &Result{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/login/%s", c.SessionKey)
	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoAndUnmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) RealtimeAccess() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/realtime_access", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) NewsSources() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/news_sources", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) NewsItems(query, sourceIds string, count, after int64) (*ResultList, error) {
	res := &ResultList{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/news_items")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	if query != "" {
		reqQuery.Set("query", query)
	}
	if sourceIds != "" {
		reqQuery.Set("sourceids", sourceIds)
	}
	if count != 0 {
		reqQuery.Set("count", strconv.FormatInt(count, 10))
	}
	if after != 0 {
		reqQuery.Set("after", strconv.FormatInt(after, 10))
	}
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) NewsItem(newsItemId int64) (*Result, error) {
	res := &Result{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/news_items/%d", newsItemId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Accounts() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/accounts", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Account(accountId int64) (*Result, error) {
	res := &Result{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/accounts/%d", accountId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) AccountLedgers(accountId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/accounts/%d/ledgers", accountId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) AccountPositions(accountId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/accounts/%d/positions", accountId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) AccountOrders(accountId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/accounts/%d/orders", accountId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) AccountTrades(accountId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/accounts/%d/trades", accountId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Instruments(query, typ, country string) (*ResultList, error) {
	res := &ResultList{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/instruments")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	if query != "" {
		reqQuery.Set("query", query)
	}
	if typ != "" {
		reqQuery.Set("type", typ)
	}
	if country != "" {
		reqQuery.Set("country", country)
	}
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Instrument(identifier string, marketId int64) (*Result, error) {
	res := &Result{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/instruments")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	reqQuery.Set("identifier", identifier)
	reqQuery.Set("marketID", strconv.FormatInt(marketId, 10))
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) ChartData(identifier string, marketId int64) (*ResultList, error) {
	res := &ResultList{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/chart_data")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	reqQuery.Set("identifier", identifier)
	reqQuery.Set("marketID", strconv.FormatInt(marketId, 10))
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Lists() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/lists", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) List(listId int64) (*ResultList, error) {
	res := &ResultList{}

	url := fmt.Sprintf("https://api.test.nordnet.se/next/v1/lists/%d", listId)
	err := c.GetAndUnmarshal(url, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Markets() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/markets", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) MarketTradingDays(marketId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/markets/%d/trading_days", marketId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Indices() (*ResultList, error) {
	res := &ResultList{}

	err := c.GetAndUnmarshal("https://api.test.nordnet.se/next/v1/indices", res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Ticksizes(ticksizeId int64) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/ticksizes/%d", ticksizeId)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) DerivateCountries(derType string) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/derivatives/%s", derType)
	err := c.GetAndUnmarshal(urlStr, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) DerivateUnderlyings(derType, country string) (*ResultList, error) {
	res := &ResultList{}

	url := fmt.Sprintf("https://api.test.nordnet.se/next/v1/derivatives/%s/underlyings/%s", derType, country)
	err := c.GetAndUnmarshal(url, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) Derivatives(derType string, marketId int64, identifier string) (*ResultList, error) {
	res := &ResultList{}

	urlStr := fmt.Sprintf("https://api.test.nordnet.se/next/v1/derivatives/%s/derivatives", derType)
	reqUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	reqQuery.Set("marketID", strconv.FormatInt(marketId, 10))
	reqQuery.Set("identifier", identifier)
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) RelatedMarkets(identifier string, marketId int64) (*ResultList, error) {
	res := &ResultList{}

	reqUrl, err := url.Parse("https://api.test.nordnet.se/next/v1/related_markets")
	if err != nil {
		return nil, err
	}

	reqQuery := reqUrl.Query()
	reqQuery.Set("identifier", identifier)
	reqQuery.Set("marketID", strconv.FormatInt(marketId, 10))
	reqUrl.RawQuery = reqQuery.Encode()

	err = c.GetAndUnmarshal(reqUrl.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) GetAndUnmarshal(url string, res interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return c.DoAndUnmarshal(req, res)
}

func (c *APIClient) DoAndUnmarshal(req *http.Request, res interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")

	if c.SessionKey != "" {
		req.SetBasicAuth(c.SessionKey, c.SessionKey)
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	c.LastUsageAt = time.Now()

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	return nil
}
