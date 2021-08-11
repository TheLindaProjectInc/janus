package transformer

import (
	"fmt"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEthValueToMetrixAmount(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in":   "0xde0b6b3a7640000",
			"want": decimal.NewFromFloat(1),
		},
		{

			"in":   "0x6f05b59d3b20000",
			"want": decimal.NewFromFloat(0.5),
		},
		{
			"in":   "0x2540be400",
			"want": decimal.NewFromFloat(0.00000001),
		},
		{
			"in":   "0x1",
			"want": decimal.NewFromInt(0),
		},
	}
	for _, c := range cases {
		in := c["in"].(string)
		want := c["want"].(decimal.Decimal)
		got, err := EthValueToMetrixAmount(in, MinimumGas)
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf("in: %s, want: %v, got: %v", in, want, got)
		}
	}
}

func TestMetrixValueToEthAmount(t *testing.T) {
	cases := []decimal.Decimal{
		decimal.NewFromFloat(1),
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.00000001),
		MinimumGas,
	}
	for _, c := range cases {
		in := c
		eth := MetrixDecimalValueToETHAmount(in)
		out := EthDecimalValueToMetrixAmount(eth)

		if !in.Equals(out) {
			t.Errorf("in: %s, eth: %v, metrix: %v", in, eth, out)
		}
	}
}

func TestMetrixAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.1), "0x16345785d8a0000"
	got, err := formatMetrixAmount(in)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("in: %v, want: %s, got: %s", in, want, got)
	}
}

func TestLowestMetrixAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.00000001), "0x2540be400"
	got, err := formatMetrixAmount(in)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("in: %v, want: %s, got: %s", in, want, got)
	}
}

func TestAddressesConversion(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		metrixChain   string
		ethAddress  string
		metrixAddress string
	}{
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "6c89a1a6ca2ae7c00b248bb2832d6f480f27da68",
			metrixAddress: "qTTH1Yr2eKCuDLqfxUyBLCAjmomQ8pyrBt",
		},

		// Test cases for addresses defined here:
		// 	- https://github.com/hayeah/openzeppelin-solidity/blob/metrixcoin/METRIX-NOTES.md#create-test-accounts
		//
		// NOTE: Ethereum addresses are without `0x` prefix, as it expects by conversion functions
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "7926223070547d2d15b2ef5e7383e541c338ffe9",
			metrixAddress: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "2352be3db3177f0a07efbe6da5857615b8c9901d",
			metrixAddress: "qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "69b004ac2b3993bf2fdf56b02746a1f57997420d",
			metrixAddress: "qTCCy8qy7pW94EApdoBjYc1vQ2w68UnXPi",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "8c647515f03daeefd09872d7530fa8d8450f069a",
			metrixAddress: "qWMi6ne9mDQFatRGejxdDYVUV9rQVkAFGp",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "2191744eb5ebeac90e523a817b77a83a0058003b",
			metrixAddress: "qLcshhsRS6HKeTKRYFdpXnGVZxw96QQcfm",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "88b0bf4b301c21f8a47be2188bad6467ad556dcf",
			metrixAddress: "qW28njWueNpBXYWj2KDmtFG2gbLeALeHfV",
		},
	}

	for i, in := range inputs {
		var (
			in       = in
			testDesc = fmt.Sprintf("#%d", i)
		)
		t.Run(testDesc, func(t *testing.T) {
			metrixAddress, err := convertETHAddress(in.ethAddress, in.metrixChain)
			require.NoError(t, err, "couldn't convert Ethereum address to Metrix address")
			require.Equal(t, in.metrixAddress, metrixAddress, "unexpected converted Metrix address value")

			ethAddress, err := convertMetrixAddress(in.metrixAddress)
			require.NoError(t, err, "couldn't convert Metrix address to Ethereum address")
			require.Equal(t, in.ethAddress, ethAddress, "unexpected converted Ethereum address value")
		})
	}
}

func TestSendTransactionRequestHasDefaultGasPriceAndAmount(t *testing.T) {
	var req eth.SendTransactionRequest
	err := unmarshalRequest([]byte(`[{}]`), &req)
	if err != nil {
		t.Fatal(err)
	}
	defaultGasPriceInWei := req.GasPrice.Int
	defaultGasPriceInMETRIX := EthDecimalValueToMetrixAmount(decimal.NewFromBigInt(defaultGasPriceInWei, 1))
	if !defaultGasPriceInMETRIX.Equals(MinimumGas) {
		t.Fatalf("Default gas price does not convert to METRIX minimum gas price, got: %s want: %s", defaultGasPriceInMETRIX.String(), MinimumGas.String())
	}
	if eth.DefaultGasAmountForMetrix.String() != req.Gas.Int.String() {
		t.Fatalf("Default gas amount does not match expected default, got: %s want: %s", req.Gas.Int.String(), eth.DefaultGasAmountForMetrix.String())
	}
}
