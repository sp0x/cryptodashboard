// Create a coingecko struct that will hold the key and the url, later on it will be used to fetch the data from the coingecko demo api.
package coingecko

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const apiHeaderName = "x-cg-demo-api-key"
const CoingeckoDefaultBaseUrl = "https://api.coingecko.com/api/v3"

type Price map[string]PriceValues
type PriceValues map[string]float64

type CoingeckoClient struct {
	httpClient *http.Client
	key        string
	baseUrl    string
}

// Create a new function that will return a new instance of the coingecko struct.
func NewCoingeckoClient(httpClient *http.Client, key string, baseUrl string) *CoingeckoClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &CoingeckoClient{
		httpClient: httpClient,
		key:        key,
		baseUrl:    baseUrl,
	}
}

func GetSimplePriceEndpoint(targetCurrency string, cryptoCurrencyIds []string) string {
	params := url.Values{}
	params.Add("ids", strings.Join(cryptoCurrencyIds, ","))
	params.Add("vs_currencies", targetCurrency)
	params.Add("include_market_cap", strconv.FormatBool(false))
	params.Add("include_24hr_vol", strconv.FormatBool(false))
	params.Add("include_24hr_change", strconv.FormatBool(false))
	params.Add("include_last_updated_at", strconv.FormatBool(false))
	params.Add("precision", "4")
	return fmt.Sprintf("/simple/price?%s", params.Encode())
}

func (c *CoingeckoClient) Get(url string) ([]byte, error) {
	fullurl := fmt.Sprintf("%s%s", c.baseUrl, url)
	req, err := http.NewRequest("GET", fullurl, nil)
	if err != nil {
		return nil, err
	}

	if c.key != "" {
		req.Header.Set(apiHeaderName, c.key)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}
