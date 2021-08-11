package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyETHUninstallFilter implements ETHProxy
type ProxyETHUninstallFilter struct {
	*metrix.Metrix
	filter *eth.FilterSimulator
}

func (p *ProxyETHUninstallFilter) Method() string {
	return "eth_uninstallFilter"
}

func (p *ProxyETHUninstallFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.UninstallFilterRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHUninstallFilter) request(ethreq *eth.UninstallFilterRequest) (eth.UninstallFilterResponse, error) {
	id, err := hexutil.DecodeUint64(string(*ethreq))
	if err != nil {
		return false, err
	}

	// uninstall
	p.filter.Uninstall(id)

	return true, nil
}
