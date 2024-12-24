package odos

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
	_baseURL = "https://api.odos.xyz"
)

type PriceResponse struct {
	CurrencyId string  `json:"currencyId"`
	Price      float64 `json:"price"`
}

type InputToken struct {
	TokenAddress string `json:"tokenAddress"`
	Amount       string `json:"amount"`
}

type OutputToken struct {
	TokenAddress string  `json:"tokenAddress"`
	Proportion   float64 `json:"proportion"`
}

type QuoteRequest struct {
	ChainId              int           `json:"chainId"`
	InputTokens          []InputToken  `json:"inputTokens"`
	OutputTokens         []OutputToken `json:"outputTokens"`
	GasPrice             float64       `json:"gasPrice"`
	UserAddr             string        `json:"userAddr"`
	SlippageLimitPercent float64       `json:"slippageLimitPercent"` // Slippage percent to use for checking if the path is valid. Float. Example: to set slippage to 0.5% send 0.5. If 1% is desired, send 1. If not provided, slippage will be set 0.3.
	SourceBlacklist      []string      `json:"sourceBlacklist"`
	SourceWhitelist      []string      `json:"sourceWhitelist"`
	PoolBlacklist        []string      `json:"poolBlacklist"`
	PathViz              bool          `json:"pathViz"`
	ReferralCode         int           `json:"referralCode"`
	Compact              bool          `json:"compact"`
	LikeAsset            bool          `json:"likeAsset"`
	DisableRFQs          bool          `json:"disableRFQs"`
	Simple               bool          `json:"simple"` // If a less complicated quote and/or a quicker response time is desired, this flag can be set. Defaults to false
}

// Token represents token information in path visualization
type Token struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Visible  bool   `json:"visible"`
	Width    int    `json:"width"`
}

// TokenInfo represents detailed token information in path links
type TokenInfo struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   int    `json:"decimals"`
	AssetID    string `json:"asset_id"`
	AssetType  string `json:"asset_type"`
	IsRebasing bool   `json:"is_rebasing"`
	CgID       string `json:"cgid"`
}

// PathLink represents a link in the path visualization
type PathLink struct {
	Source       int       `json:"source"`
	Target       int       `json:"target"`
	SourceExtend bool      `json:"sourceExtend"`
	TargetExtend bool      `json:"targetExtend"`
	Label        string    `json:"label"`
	Value        float64   `json:"value"`
	NextValue    float64   `json:"nextValue"`
	StepValue    float64   `json:"stepValue"`
	InValue      float64   `json:"in_value"`
	OutValue     float64   `json:"out_value"`
	EdgeLen      int       `json:"edge_len"`
	SourceToken  TokenInfo `json:"sourceToken"`
	TargetToken  TokenInfo `json:"targetToken"`
}

// PathViz represents the path visualization
type PathViz struct {
	Nodes []Token    `json:"nodes"`
	Links []PathLink `json:"links"`
}

// QuoteResponse represents the response from quote endpoint
type QuoteResponse struct {
	InTokens          []string  `json:"inTokens"`
	OutTokens         []string  `json:"outTokens"`
	InAmounts         []string  `json:"inAmounts"`
	OutAmounts        []string  `json:"outAmounts"`
	GasEstimate       float64   `json:"gasEstimate"`
	DataGasEstimate   int       `json:"dataGasEstimate"`
	GweiPerGas        float64   `json:"gweiPerGas"`
	GasEstimateValue  float64   `json:"gasEstimateValue"`
	InValues          []float64 `json:"inValues"`
	OutValues         []float64 `json:"outValues"`
	NetOutValue       float64   `json:"netOutValue"`
	PriceImpact       float64   `json:"priceImpact"`
	PercentDiff       float64   `json:"percentDiff"`
	PartnerFeePercent float64   `json:"partnerFeePercent"`
	PathId            string    `json:"pathId"`
	PathViz           PathViz   `json:"pathViz"`
	BlockNumber       int64     `json:"blockNumber"`
}

// AssembleRequest represents the request body for assemble endpoint
type AssembleRequest struct {
	UserAddr string `json:"userAddr"`
	PathId   string `json:"pathId"`
	Simulate bool   `json:"simulate"`
}

// Transaction represents the transaction details in the assemble response
type Transaction struct {
	Gas      int64  `json:"gas"`
	GasPrice int64  `json:"gasPrice"`
	Value    string `json:"value"`
	To       string `json:"to"`
	From     string `json:"from"`
	Data     string `json:"data"`
	Nonce    int64  `json:"nonce"`
	ChainId  int    `json:"chainId"`
}

// Simulation represents the simulation results
type Simulation struct {
	IsSuccess       bool    `json:"isSuccess"`
	AmountsOut      []int64 `json:"amountsOut"`
	GasEstimate     int64   `json:"gasEstimate"`
	SimulationError string  `json:"simulationError"`
}

// AssembleResponse represents the response from assemble endpoint
type AssembleResponse struct {
	Deprecated       *string      `json:"deprecated"`
	BlockNumber      int64        `json:"blockNumber"`
	GasEstimate      int64        `json:"gasEstimate"`
	GasEstimateValue float64      `json:"gasEstimateValue"`
	InputTokens      []InputToken `json:"inputTokens"`
	OutputTokens     []struct {
		TokenAddress string `json:"tokenAddress"`
		Amount       string `json:"amount"`
	} `json:"outputTokens"`
	NetOutValue float64     `json:"netOutValue"`
	OutValues   []string    `json:"outValues"`
	Transaction Transaction `json:"transaction"`
	Simulation  Simulation  `json:"simulation"`
}

type OdosClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new KyberSwap client
func NewClient(baseURL string) *OdosClient {
	if baseURL == "" {
		baseURL = _baseURL
	}

	return &OdosClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *OdosClient) GetTokenPrice(chainID, tokenAddr string) (*PriceResponse, error) {
	url := fmt.Sprintf("%s/pricing/token/%s/%s", c.baseURL, chainID, tokenAddr)
	log.Info().Msgf("url: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get token price: %w", err)
	}
	defer resp.Body.Close()

	var priceResp PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &priceResp, nil
}

// Generate Odos Quote
// /sor/quote/v2
func (c *OdosClient) Quote(req *QuoteRequest) (*QuoteResponse, error) {
	url := fmt.Sprintf("%s/sor/quote/v2", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Origin", "https://app.odos.xyz")
	request.Header.Set("Referer", "https://app.odos.xyz/")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}
	defer resp.Body.Close()

	var quoteResp QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &quoteResp, nil
}

// /sor/assemble
// Assemble Odos quote into transaction
func (c *OdosClient) Assemble(userAddr, pathId string, isSimulate bool) (*AssembleResponse, error) {
	url := fmt.Sprintf("%s/sor/assemble", c.baseURL)

	req := AssembleRequest{
		UserAddr: userAddr,
		PathId:   pathId,
		Simulate: isSimulate,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Origin", "https://app.odos.xyz")
	request.Header.Set("Referer", "https://app.odos.xyz/")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Info().Msgf("response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response_body", string(body)).
			Msg("Assemble request failed")
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var assembleResp AssembleResponse
	if err := json.Unmarshal(body, &assembleResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &assembleResp, nil
}
