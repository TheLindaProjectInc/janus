package transformer

import (
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

// ProxyETHGetCode implements ETHProxy
type ProxyETHGetCode struct {
	*metrix.Metrix
}

func (p *ProxyETHGetCode) Method() string {
	return "eth_getCode"
}

func (p *ProxyETHGetCode) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.GetCodeRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHGetCode) request(ethreq *eth.GetCodeRequest) (eth.GetCodeResponse, error) {
	metrixreq := metrix.GetAccountInfoRequest(utils.RemoveHexPrefix(ethreq.Address))

	metrixresp, err := p.GetAccountInfo(&metrixreq)
	if err != nil {
		if err == metrix.ErrInvalidAddress {
			/**
			// correct response for an invalid address
			{
				"jsonrpc": "2.0",
				"id": 123,
				"result": "0x"
			}
			**/
			return "0x", nil
		} else {
			return "", err
		}
	}

	// metrix res -> eth res
	return eth.GetCodeResponse(utils.AddHexPrefix(metrixresp.Code)), nil
}
