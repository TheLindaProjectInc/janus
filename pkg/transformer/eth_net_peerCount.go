package transformer

import (
	"github.com/dcb9/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyNetPeerCount implements ETHProxy
type ProxyNetPeerCount struct {
	*metrix.Metrix
}

func (p *ProxyNetPeerCount) Method() string {
	return "net_peerCount"
}

func (p *ProxyNetPeerCount) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyNetPeerCount) request() (*eth.NetPeerCountResponse, error) {
	peerInfos, err := p.GetPeerInfo()
	if err != nil {
		return nil, err
	}

	resp := eth.NetPeerCountResponse(hexutil.EncodeUint64(uint64(len(peerInfos))))
	return &resp, nil
}
