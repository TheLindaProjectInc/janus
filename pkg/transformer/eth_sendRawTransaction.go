package transformer

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

// ProxyETHSendRawTransaction implements ETHProxy
type ProxyETHSendRawTransaction struct {
	*metrix.Metrix
}

var _ ETHProxy = (*ProxyETHSendRawTransaction)(nil)

func (p *ProxyETHSendRawTransaction) Method() string {
	return "eth_sendRawTransaction"
}

func (p *ProxyETHSendRawTransaction) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var params eth.SendRawTransactionRequest
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, err
	}
	if params[0] == "" {
		return nil, errors.Errorf("invalid parameter: raw transaction hexed string is empty")
	}

	return p.request(params)
}

func (p *ProxyETHSendRawTransaction) request(params eth.SendRawTransactionRequest) (eth.SendRawTransactionResponse, error) {
	var (
		metrixHexedRawTx = utils.RemoveHexPrefix(params[0])
		req            = metrix.SendRawTransactionRequest([1]string{metrixHexedRawTx})
	)

	metrixresp, err := p.Metrix.SendRawTransaction(&req)
	if err != nil {
		if err == metrix.ErrVerifyAlreadyInChain {
			// already committed
			// we need to send back the tx hash
			rawTx, err := p.Metrix.DecodeRawTransaction(metrixHexedRawTx)
			if err != nil {
				p.GetErrorLogger().Log("msg", "Error decoding raw transaction for duplicate raw transaction", "err", err)
				return eth.SendRawTransactionResponse(""), err
			}
			metrixresp = &metrix.SendRawTransactionResponse{Result: rawTx.Hash}
		} else {
			return eth.SendRawTransactionResponse(""), err
		}
	} else {
		if p.Chain() == metrix.ChainRegTest {
			if _, err = p.Generate(1, nil); err != nil {
				p.GetErrorLogger().Log("Error generating new block", err)
			}
		}
	}

	resp := *metrixresp
	ethHexedTxHash := utils.AddHexPrefix(resp.Result)
	return eth.SendRawTransactionResponse(ethHexedTxHash), nil
}
