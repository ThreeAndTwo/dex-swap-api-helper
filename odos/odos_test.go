package odos

import (
	"encoding/json"
	"testing"
)

const (
	chainId = "1"
	sUSDe   = "0x9D39A5DE30e57443BfF2A8307A4256c8797A3497"
	DAI     = "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	wstETH = "0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0"
	ezETH  = "0xbf5495Efe5DB9ce00f80364C8B423567e58d2110"
)

var odosClient *OdosClient

func init() {
	odosClient = NewClient("") // baseURL is empty, so it will use the default baseURL
}

func TestGetTokenPrice(t *testing.T) {
	type args struct {
		chainID   string
		tokenAddr string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test get router by DAI",
			args: args{
				chainID:   chainId,
				tokenAddr: DAI,
			},
			wantErr: false,
		},
		{
			name: "test get router by USDC",
			args: args{
				chainID:   chainId,
				tokenAddr: sUSDe,
			},
			wantErr: false,
		},
		{
			name: "test get router by wstETH",
			args: args{
				chainID:   chainId,
				tokenAddr: wstETH,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := odosClient.GetTokenPrice(tt.args.chainID, tt.args.tokenAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenPrice() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("GetTokenPrice() = %v", got)
		})
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		name    string
		args    *QuoteRequest
		wantErr bool
	}{
		{
			name: "test get router by DAI",
			args: &QuoteRequest{
				ChainId: 1,
				InputTokens: []InputToken{
					{
						TokenAddress: DAI,
						Amount:       "1000000000000000000",
					},
				},
				OutputTokens: []OutputToken{
					{
						TokenAddress: sUSDe,
						Proportion:   1,
					},
				},
				GasPrice:             6.27,
				UserAddr:             "0x0000000000000000000000000000000000000000",
				SlippageLimitPercent: 0.1, // 0.1%
				SourceBlacklist:      []string{},
				SourceWhitelist:      []string{},
				PoolBlacklist:        []string{},
				PathViz:              true,
				ReferralCode:         1,
				Compact:              true,
				LikeAsset:            true,
				DisableRFQs:          false,
				Simple:               false,
			},
			wantErr: false,
		},
		{
			name: "test get router by USDC",
			args: &QuoteRequest{
				ChainId: 1,
				InputTokens: []InputToken{
					{
						TokenAddress: sUSDe,
						Amount:       "1000000000000000000",
					},
				},
				OutputTokens: []OutputToken{
					{
						TokenAddress: DAI,
						Proportion:   1,
					},
				},
				GasPrice:             6.27,
				UserAddr:             "0x0000000000000000000000000000000000000000",
				SlippageLimitPercent: 0.1, // 0.1%
				SourceBlacklist:      []string{},
				SourceWhitelist:      []string{},
				PoolBlacklist:        []string{},
				PathViz:              true,
				ReferralCode:         1,
				Compact:              true,
				LikeAsset:            true,
				DisableRFQs:          false,
				Simple:               false,
			},
			wantErr: false,
		},
		{
			name: "test get router by wstETH",
			args: &QuoteRequest{
				ChainId: 1,
				InputTokens: []InputToken{
					{
						TokenAddress: wstETH,
						Amount:       "1000000000000000000",
					},
				},
				OutputTokens: []OutputToken{
					{
						TokenAddress: ezETH,
						Proportion:   1,
					},
				},
				GasPrice:             6.27,
				UserAddr:             "0x0000000000000000000000000000000000000000",
				SlippageLimitPercent: 0.1, // 0.1%
				SourceBlacklist:      []string{},
				SourceWhitelist:      []string{},
				PoolBlacklist:        []string{},
				PathViz:              true,
				ReferralCode:         1,
				Compact:              true,
				LikeAsset:            true,
				DisableRFQs:          false,
				Simple:               false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := odosClient.Quote(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quote() error = %v, wantErr %v", err, tt.wantErr)
			}

			marData, _ := json.Marshal(got)
			t.Logf("Quote() = %v", string(marData))
		})
	}
}

func TestAssemble(t *testing.T) {
	type args struct {
		userAddr string
		pathId   string
		simulate bool
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test assemble ETH -> sUSDe",
			args: args{
				userAddr: "0x163A5EC5e9C32238d075E2D829fE9fA87451e3b7",
				pathId:   "9c2294c5e076d888e149c764f832738b",
				simulate: true,
			},
			wantErr: false,
		},
		{
			name: "test assemble DAI -> sUSDe",
			args: args{
				userAddr: "0x163A5EC5e9C32238d075E2D829fE9fA87451e3b7",
				pathId:   "5336207d5a6c0bd9286671ba4640aa0d",
				simulate: true,
			},
			wantErr: false,
		},
		// {
		// 	name: "test assemble sUSDe -> DAI",
		// 	args: args{
		// 		userAddr: "0x0000000000000000000000000000000000000000",
		// 		pathId:   "d257e17a73104028e35630a6ba6b7952",
		// 		simulate: true,
		// 	},
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := odosClient.Assemble(tt.args.userAddr, tt.args.pathId, tt.args.simulate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			marData, _ := json.Marshal(got)
			t.Logf("Assemble() = %v", string(marData))
		})
	}
}
