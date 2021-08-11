package transformer

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

// ProxyETHSendTransaction implements ETHProxy
type ProxyETHSendTransaction struct {
	*metrix.Metrix
}

func (p *ProxyETHSendTransaction) Method() string {
	return "eth_sendTransaction"
}

func (p *ProxyETHSendTransaction) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.SendTransactionRequest
	err := unmarshalRequest(rawreq.Params, &req)
	if err != nil {
		return nil, err
	}

	var result interface{}

	if req.IsCreateContract() {
		result, err = p.requestCreateContract(&req)
	} else if req.IsSendEther() {
		result, err = p.requestSendToAddress(&req)
	} else if req.IsCallContract() {
		result, err = p.requestSendToContract(&req)
	} else {
		return nil, errors.New("Unknown operation")
	}

	if p.Chain() == metrix.ChainRegTest && err == nil {
		defer func() {
			if _, generateErr := p.Generate(1, nil); generateErr != nil {
				p.GetErrorLogger().Log("Error generating new block", generateErr)
			}
		}()
	}

	return result, err
}

func (p *ProxyETHSendTransaction) requestSendToContract(ethtx *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	gasLimit, gasPrice, err := EthGasToMetrix(ethtx)
	if err != nil {
		return nil, err
	}

	amount := decimal.NewFromFloat(0.0)
	if ethtx.Value != "" {
		var err error
		amount, err = EthValueToMetrixAmount(ethtx.Value, ZeroSatoshi)
		if err != nil {
			return nil, errors.Wrap(err, "EthValueToMetrixAmount:")
		}
	}

	metrixreq := metrix.SendToContractRequest{
		ContractAddress: utils.RemoveHexPrefix(ethtx.To),
		Datahex:         utils.RemoveHexPrefix(ethtx.Data),
		Amount:          amount,
		GasLimit:        gasLimit,
		GasPrice:        gasPrice,
	}

	if from := ethtx.From; from != "" && utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
		metrixreq.SenderAddress = from
	}

	var resp *metrix.SendToContractResponse
	if err := p.Metrix.Request(metrix.MethodSendToContract, &metrixreq, &resp); err != nil {
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(utils.AddHexPrefix(resp.Txid))
	return &ethresp, nil
}

func (p *ProxyETHSendTransaction) requestSendToAddress(req *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	getMetrixWalletAddress := func(addr string) (string, error) {
		if utils.IsEthHexAddress(addr) {
			return p.FromHexAddress(utils.RemoveHexPrefix(addr))
		}
		return addr, nil
	}

	from, err := getMetrixWalletAddress(req.From)
	if err != nil {
		return nil, err
	}

	to, err := getMetrixWalletAddress(req.To)
	if err != nil {
		return nil, err
	}

	amount, err := EthValueToMetrixAmount(req.Value, ZeroSatoshi)
	if err != nil {
		return nil, errors.Wrap(err, "EthValueToMetrixAmount:")
	}

	p.GetDebugLogger().Log("msg", "successfully converted from wei to MRX", "wei", req.Value, "metrix", amount)

	metrixreq := metrix.SendToAddressRequest{
		Address:       to,
		Amount:        amount,
		SenderAddress: from,
	}

	var metrixresp metrix.SendToAddressResponse
	if err := p.Metrix.Request(metrix.MethodSendToAddress, &metrixreq, &metrixresp); err != nil {
		// this can fail with:
		// "error": {
		//   "code": -3,
		//   "message": "Sender address does not have any unspent outputs"
		// }
		// this can happen if there are enough coins but some required are untrusted
		// you can get the trusted coin balance via getbalances rpc call
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(utils.AddHexPrefix(string(metrixresp)))

	return &ethresp, nil
}

func (p *ProxyETHSendTransaction) requestCreateContract(req *eth.SendTransactionRequest) (*eth.SendTransactionResponse, error) {
	gasLimit, gasPrice, err := EthGasToMetrix(req)
	if err != nil {
		return nil, err
	}

	metrixreq := &metrix.CreateContractRequest{
		ByteCode: utils.RemoveHexPrefix(req.Data),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	if req.From != "" {
		from := req.From
		if utils.IsEthHexAddress(from) {
			from, err = p.FromHexAddress(from)
			if err != nil {
				return nil, err
			}
		}

		metrixreq.SenderAddress = from
	}

	var resp *metrix.CreateContractResponse
	if err := p.Metrix.Request(metrix.MethodCreateContract, metrixreq, &resp); err != nil {
		return nil, err
	}

	ethresp := eth.SendTransactionResponse(utils.AddHexPrefix(string(resp.Txid)))

	return &ethresp, nil
}
