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
			metrixAddress: "mS5FATs9YWa931SoiHdzUwSsUmBy28nHrr",
		},

		// Test cases for addresses defined here:
		// 	- https://github.com/hayeah/openzeppelin-solidity/blob/qtum/METRIX-NOTES.md#create-test-accounts
		//
		// NOTE: Ethereum addresses are without `0x` prefix, as it expects by conversion functions
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "b89a4201258da334e3cd6d49047715fbf8a0e386",
			metrixAddress: "mZ1SSGGtAav5b5rCgf4x5SphNLLv6EVtMT",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "f3fdb4c8636a0fc4a6b00922b5c646650821deab",
			metrixAddress: "meRTSNhCRNDmdNkvMumNVvEiZAresrzVbV",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "4f9fd8b40f145cbd2ada7e8fcd81034e3f77d05d",
			metrixAddress: "mPSNBRriZPyKRXPwXDPovGtx9pgFm4Erjs",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "965a1bf254099bff1467e4020781e0dafe073810",
			metrixAddress: "mVtLccbttd4vCZCrf7gKCyZoeZWDjQvrS7",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "43fe2890bb0acaa1f8a56674e43d618ce425e516",
			metrixAddress: "mNNs437u1qc6DVUbeGGpY7HsJPDRB8JuB4",
		},
		{
			metrixChain:   metrix.ChainTest,
			ethAddress:  "43f7324ad7e65e30f6ac9e2c00c2f020c0f53dc0",
			metrixAddress: "mNNiiJsBnPUZu4NwNtPaK9zyxD5ghV1By8",
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
