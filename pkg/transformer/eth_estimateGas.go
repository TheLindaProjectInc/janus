package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

var ErrExecutionReverted = errors.New("execution reverted")

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHEstimateGas struct {
	*ProxyETHCall
}

func (p *ProxyETHEstimateGas) Method() string {
	return "eth_estimateGas"
}

func (p *ProxyETHEstimateGas) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var ethreq eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &ethreq); err != nil {
		return nil, err
	}

	// when supplying this parameter to callcontract to estimate gas in the metrix api
	// if there isn't enough gas specified here, the result will be an exception
	// Excepted = "OutOfGasIntrinsic"
	// Gas = "the supplied value"
	// this is different from geth's behavior
	// which will return a used gas value that is higher than the incoming gas parameter
	// so we set this to nil so that callcontract will return the actual gas estimate
	ethreq.Gas = nil

	// eth req -> metrix req
	metrixreq, err := p.ToRequest(&ethreq)
	if err != nil {
		return nil, err
	}

	// metrix [code: -5] Incorrect address occurs here
	metrixresp, err := p.CallContract(metrixreq)
	if err != nil {
		return nil, err
	}

	return p.toResp(metrixresp)
}

func (p *ProxyETHEstimateGas) toResp(metrixresp *metrix.CallContractResponse) (*eth.EstimateGasResponse, error) {
	if metrixresp.ExecutionResult.Excepted != "None" {
		// TODO: Return code -32000
		return nil, ErrExecutionReverted
	}
	gas := eth.EstimateGasResponse(hexutil.EncodeUint64(uint64(metrixresp.ExecutionResult.GasUsed)))
	p.GetDebugLogger().Log(p.Method(), gas)
	return &gas, nil
}
