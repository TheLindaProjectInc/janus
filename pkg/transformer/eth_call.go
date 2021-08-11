package transformer

import (
	"math/big"

	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*metrix.Metrix
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHCall) request(ethreq *eth.CallRequest) (interface{}, error) {
	// eth req -> metrix req
	metrixreq, err := p.ToRequest(ethreq)
	if err != nil {
		return nil, err
	}

	metrixresp, err := p.CallContract(metrixreq)
	if err != nil {
		return nil, err
	}

	// metrix res -> eth res
	return p.ToResponse(metrixresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*metrix.CallContractRequest, error) {
	from := ethreq.From
	var err error
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
	}

	var gasLimit *big.Int
	if ethreq.Gas != nil {
		gasLimit = ethreq.Gas.Int
	}

	return &metrix.CallContractRequest{
		To:       ethreq.To,
		From:     from,
		Data:     ethreq.Data,
		GasLimit: gasLimit,
	}, nil
}

func (p *ProxyETHCall) ToResponse(qresp *metrix.CallContractResponse) interface{} {

	if qresp.ExecutionResult.Output == "" {

		return &eth.JSONRPCError{
			Message: "Revert: executionResult output is empty",
			Code:    -32000,
		}

	}

	data := utils.AddHexPrefix(qresp.ExecutionResult.Output)
	metrixresp := eth.CallResponse(data)
	return &metrixresp

}
