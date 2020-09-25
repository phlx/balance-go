package exchangerates

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	baseUrl = "https://api.exchangeratesapi.io"
)

type LatestResponse struct {
	Rates map[Currency]float64 `json:"rates"`
	Base  Currency             `json:"base"`
	Date  string               `json:"string"`
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

type Client struct {
	httpClient *http.Client
}

func (c *Client) Latest(ctx context.Context) (*LatestResponse, error) {
	url := baseUrl + "/latest?base=" + string(RUB)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var latestResponse LatestResponse
	err = json.Unmarshal(body, &latestResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON from string '%s'", body)
	}
	return &latestResponse, nil
}
