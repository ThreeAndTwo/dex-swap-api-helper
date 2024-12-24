package kyberswap

import (
	"encoding/json"
	"testing"
)

//https://aggregator-api.kyberswap.com/ethereum/api/v1/routes?tokenIn=0x9D39A5DE30e57443BfF2A8307A4256c8797A3497&tokenOut=0xdC035D45d973E3EC169d2276DDab16f1e407384F&amountIn=2000000000000000000000000&gasInclude=true

const (
	chain = "ethereum"
	sUSDe = "0x9D39A5DE30e57443BfF2A8307A4256c8797A3497"

	USDT = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	USDC = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	DAI  = "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	wstETH = "0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0"
	ezETH  = "0xbf5495Efe5DB9ce00f80364C8B423567e58d2110"
)

var kyberSwapClient *KyberSwapClient

func init() {
	kyberSwapClient = NewClient("", chain) // baseURL is empty, so it will use the default baseURL
}

func TestKyberSwapClient_GetRoutes(t *testing.T) {
	type args struct {
		tokenIn  string
		tokenOut string
		amountIn string
	}
	tests := []struct {
		name    string
		args    args
		want    *RouteResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test get router by USDT",
			args: args{
				tokenIn:  USDT,
				tokenOut: sUSDe,
				amountIn: "2238451.467827",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test get router by USDC",
			args: args{
				tokenIn:  DAI,
				tokenOut: sUSDe,
				amountIn: "100",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test get router by wstETH",
			args: args{
				tokenIn:  wstETH,
				tokenOut: ezETH,
				amountIn: "10000000000000",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kyberSwapClient.GetRoutes(tt.args.tokenIn, tt.args.tokenOut, tt.args.amountIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			marshal, _ := json.Marshal(got)
			t.Log(string(marshal))
			t.Log("********************************************************")
			sender := "0xd46B96d15ffF9b2B17e9c788086f3159bD0e8355"

			route, err := kyberSwapClient.BuildRoute(got.Data.RouteSummary, sender, sender)
			if err != nil {
				t.Errorf("kyberSwapClient.GetRoutes() error = %v", err)
				return
			}
			marshal, _ = json.Marshal(route)
			t.Log(string(marshal))

			t.Log("--------------------------------------------------------")
		})
	}
}
