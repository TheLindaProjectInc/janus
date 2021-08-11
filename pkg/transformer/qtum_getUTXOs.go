package transformer

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

type ProxyMETRIXGetUTXOs struct {
	*metrix.Metrix
}

var _ ETHProxy = (*ProxyMETRIXGetUTXOs)(nil)

func (p *ProxyMETRIXGetUTXOs) Method() string {
	return "metrix_getUTXOs"
}

func (p *ProxyMETRIXGetUTXOs) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var params eth.GetUTXOsRequest
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, errors.WithMessage(err, "couldn't unmarshal request parameters")
	}

	err := params.CheckHasValidValues()
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't validate parameters value")
	}

	return p.request(params)
}

func (p *ProxyMETRIXGetUTXOs) request(params eth.GetUTXOsRequest) (*eth.GetUTXOsResponse, error) {
	address, err := convertETHAddress(utils.RemoveHexPrefix(params.Address), p.Chain())
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't convert Ethereum address to Metrix address")
	}

	req := metrix.GetAddressUTXOsRequest{
		Addresses: []string{address},
	}

	resp, err := p.Metrix.GetAddressUTXOs(&req)
	if err != nil {
		return nil, err
	}

	//Convert minSumAmount to Satoshis
	minimumSum := convertFromMetrixToSatoshis(params.MinSumAmount)

	var utxos []eth.MetrixUTXO
	var minUTXOsSum decimal.Decimal
	for _, utxo := range *resp {
		minUTXOsSum = minUTXOsSum.Add(utxo.Satoshis)
		utxos = append(utxos, toEthResponseType(utxo))
		if minUTXOsSum.GreaterThanOrEqual(minimumSum) {
			return (*eth.GetUTXOsResponse)(&utxos), nil
		}
	}

	return nil, errors.Errorf("required minimum amount is greater than total amount of UTXOs")
}

func toEthResponseType(utxo metrix.UTXO) eth.MetrixUTXO {
	return eth.MetrixUTXO{
		Address: utxo.Address,
		TXID:    utxo.TXID,
		Vout:    utxo.OutputIndex,
		Amount:  convertFromSatoshisToMetrix(utxo.Satoshis).String(),
	}
}
