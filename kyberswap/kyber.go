package kyberswap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	_baseURL = "https://aggregator-api.kyberswap.com"
)

// Client represents a KyberSwap API client
type KyberSwapClient struct {
	httpClient *http.Client
	baseURL    string
}

// RouteResponse represents the API response structure
type RouteResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RouteSummary  RouteSummary `json:"routeSummary"`
		RouterAddress string       `json:"routerAddress"`
	} `json:"data"`
	RequestId string `json:"requestId"`
}

// RouteSummary represents the route summary information
type RouteSummary struct {
	TokenIn                      string    `json:"tokenIn"`
	AmountIn                     string    `json:"amountIn"`
	AmountInUsd                  string    `json:"amountInUsd"`
	TokenInMarketPriceAvailable  bool      `json:"tokenInMarketPriceAvailable"`
	TokenOut                     string    `json:"tokenOut"`
	AmountOut                    string    `json:"amountOut"`
	AmountOutUsd                 string    `json:"amountOutUsd"`
	TokenOutMarketPriceAvailable bool      `json:"tokenOutMarketPriceAvailable"`
	Gas                          string    `json:"gas"`
	GasPrice                     string    `json:"gasPrice"`
	GasUsd                       string    `json:"gasUsd"`
	ExtraFee                     ExtraFee  `json:"extraFee"`
	Route                        [][]Route `json:"route"`
}

// ExtraFee represents the fee information
type ExtraFee struct {
	FeeAmount   string `json:"feeAmount"`
	ChargeFeeBy string `json:"chargeFeeBy"`
	IsInBps     bool   `json:"isInBps"`
	FeeReceiver string `json:"feeReceiver"`
}

// Route represents a single route segment
type Route struct {
	Pool              string         `json:"pool"`
	TokenIn           string         `json:"tokenIn"`
	TokenOut          string         `json:"tokenOut"`
	LimitReturnAmount string         `json:"limitReturnAmount"`
	SwapAmount        string         `json:"swapAmount"`
	AmountOut         string         `json:"amountOut"`
	Exchange          string         `json:"exchange"`
	PoolLength        int            `json:"poolLength"`
	PoolType          string         `json:"poolType"`
	PoolExtra         OuterPoolExtra `json:"poolExtra"`
	Extra             interface{}    `json:"extra"`
}

type OuterPoolExtra struct {
	BlockNumber      int64 `json:"blockNumber"`
	TokenInIndex     int   `json:"tokenInIndex"`
	TokenOutIndex    int   `json:"tokenOutIndex"`
	Underlying       bool  `json:"underlying"`
	TokenInIsNative  bool  `json:"TokenInIsNative"`
	TokenOutIsNative bool  `json:"TokenOutIsNative"`
}

type BuildRouteRequest struct {
	RouteSummary      RouteSummary `json:"routeSummary"`
	Sender            string       `json:"sender"`
	Recipient         string       `json:"recipient"`
	Deadline          int64        `json:"deadline"`
	SlippageTolerance int64        `json:"slippageTolerance"`
}

// BuildRouteResponse represents the response from building a route
type BuildRouteResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AmountIn         string       `json:"amountIn"`
		AmountInUsd      string       `json:"amountInUsd"`
		AmountOut        string       `json:"amountOut"`
		AmountOutUsd     string       `json:"amountOutUsd"`
		Gas              string       `json:"gas"`
		GasUsd           string       `json:"gasUsd"`
		OutputChange     OutputChange `json:"outputChange"`
		Data             string       `json:"data"`
		RouterAddress    string       `json:"routerAddress"`
		TransactionValue string       `json:"transactionValue"`
	} `json:"data"`
	RequestId string `json:"requestId"`
}

// OutputChange represents the change in output amount
type OutputChange struct {
	Amount  string  `json:"amount"`
	Percent float64 `json:"percent"`
	Level   int     `json:"level"`
}

// NewClient creates a new KyberSwap client
func NewClient(baseURL, chain string) *KyberSwapClient {
	if baseURL == "" {
		baseURL = _baseURL
	}

	if chain == "" {
		chain = "ethereum"
	}

	return &KyberSwapClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: fmt.Sprintf("%s/%s", baseURL, chain),
	}
}

// GetRoutes fetches routes for token swap
func (c *KyberSwapClient) GetRoutes(tokenIn, tokenOut, amountIn string) (*RouteResponse, error) {
	url := fmt.Sprintf("%s/api/v1/routes?tokenIn=%s&tokenOut=%s&amountIn=%s",
		c.baseURL, tokenIn, tokenOut, amountIn)
	log.Info().Msgf("url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var routeResp RouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&routeResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &routeResp, nil
}

// BuildRoute sends a request to build a route
func (c *KyberSwapClient) BuildRoute(routeSummary RouteSummary, sender, recipient string) (*BuildRouteResponse, error) {
	reqBody := BuildRouteRequest{
		RouteSummary:      routeSummary,
		Sender:            sender,
		Recipient:         recipient,
		Deadline:          time.Now().Unix() + 20*3600, // TODO: need deleted
		SlippageTolerance: 10,                          // 0.1%
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	log.Debug().Msgf("jsonBody: %s", string(jsonBody))

	url := fmt.Sprintf("%s/api/v1/route/build", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("status code %d, failed to read error response: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, string(body))
	}

	var buildResp BuildRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &buildResp, nil
}

// WithTimeout sets a custom timeout for the HTTP client
func (c *KyberSwapClient) WithTimeout(timeout time.Duration) *KyberSwapClient {
	c.httpClient.Timeout = timeout
	return c
}
