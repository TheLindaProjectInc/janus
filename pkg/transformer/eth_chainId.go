package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

type ProxyETHChainId struct {
	*metrix.Metrix
}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var metrixresp *metrix.GetBlockChainInfoResponse
	if err := p.Metrix.Request(metrix.MethodGetBlockChainInfo, nil, &metrixresp); err != nil {
		return nil, err
	}

	var chainId *big.Int
	switch metrixresp.Chain {
	case "regtest":
		chainId = big.NewInt(113)
	default:
		chainId = big.NewInt(81)
		p.GetDebugLogger().Log("method", p.Method(), "msg", "Unknown chain "+metrixresp.Chain)
	}

	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}
