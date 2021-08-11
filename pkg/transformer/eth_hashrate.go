package transformer

import (
	"math"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

//ProxyETHGetHashrate implements ETHProxy
type ProxyETHHashrate struct {
	*metrix.Metrix
}

func (p *ProxyETHHashrate) Method() string {
	return "eth_hashrate"
}

func (p *ProxyETHHashrate) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHHashrate) request() (*eth.HashrateResponse, error) {
	metrixresp, err := p.Metrix.GetHashrate()
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.ToResponse(metrixresp), nil
}

func (p *ProxyETHHashrate) ToResponse(metrixresp *metrix.GetHashrateResponse) *eth.HashrateResponse {
	hexVal := hexutil.EncodeUint64(math.Float64bits(metrixresp.Difficulty))
	ethresp := eth.HashrateResponse(hexVal)
	return &ethresp
}
