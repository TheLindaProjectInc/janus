package transformer

import (
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

//ProxyETHGetHashrate implements ETHProxy
type ProxyETHMining struct {
	*metrix.Metrix
}

func (p *ProxyETHMining) Method() string {
	return "eth_mining"
}

func (p *ProxyETHMining) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHMining) request() (*eth.MiningResponse, error) {
	metrixresp, err := p.Metrix.GetMining()
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.ToResponse(metrixresp), nil
}

func (p *ProxyETHMining) ToResponse(metrixresp *metrix.GetMiningResponse) *eth.MiningResponse {
	ethresp := eth.MiningResponse(metrixresp.Staking)
	return &ethresp
}
