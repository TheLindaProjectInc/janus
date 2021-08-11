package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHGasPrice struct {
	*metrix.Metrix
}

func (p *ProxyETHGasPrice) Method() string {
	return "eth_gasPrice"
}

func (p *ProxyETHGasPrice) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	metrixresp, err := p.Metrix.GetGasPrice()
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.response(metrixresp), nil
}

func (p *ProxyETHGasPrice) response(metrixresp *big.Int) string {
	return hexutil.EncodeBig(metrixresp)
}
