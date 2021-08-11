package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHTxCount struct {
	*metrix.Metrix
}

func (p *ProxyETHTxCount) Method() string {
	return "eth_getTransactionCount"
}

func (p *ProxyETHTxCount) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {

	/* not sure we need this. Need to figure out how to best unmarshal this in the future. For now this will work.
	var req eth.GetTransactionCountRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}*/
	metrixresp, err := p.Metrix.GetTransactionCount("", "")
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.response(metrixresp), nil
}

func (p *ProxyETHTxCount) response(metrixresp *big.Int) string {
	return hexutil.EncodeBig(metrixresp)
}
