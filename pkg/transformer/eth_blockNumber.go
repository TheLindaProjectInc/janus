package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*metrix.Metrix
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHBlockNumber) request() (*eth.BlockNumberResponse, error) {
	metrixresp, err := p.Metrix.GetBlockCount()
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.ToResponse(metrixresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(metrixresp *metrix.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(metrixresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
