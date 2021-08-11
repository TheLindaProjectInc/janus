package transformer

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

// ProxyETHSendTransaction implements ETHProxy
type ProxyETHSignTransaction struct {
	*metrix.Metrix
}

func (p *ProxyETHSignTransaction) Method() string {
	return "eth_signTransaction"
}

func (p *ProxyETHSignTransaction) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.SendTransactionRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	if req.IsCreateContract() {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "transaction is a create contract request")
		return p.requestCreateContract(&req)
	} else if req.IsSendEther() {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "transaction is a send ether request")
		return p.requestSendToAddress(&req)
	} else if req.IsCallContract() {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "transaction is a call contract request")
		return p.requestSendToContract(&req)
	} else {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "transaction is an unknown request")
	}

	return nil, errors.New("Unknown operation")
}

func (p *ProxyETHSignTransaction) getRequiredUtxos(from string, neededAmount decimal.Decimal) ([]metrix.RawTxInputs, decimal.Decimal, error) {
	//convert address to metrix address
	addr := utils.RemoveHexPrefix(from)
	base58Addr, err := p.FromHexAddress(addr)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}
	// need to get utxos with txid and vouts. In order to do this we get a list of unspent transactions and begin summing them up
	var getaddressutxos *metrix.GetAddressUTXOsRequest = &metrix.GetAddressUTXOsRequest{Addresses: []string{base58Addr}}
	metrixresp, err := p.GetAddressUTXOs(getaddressutxos)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	//Convert minSumAmount to Satoshis
	minimumSum := convertFromMetrixToSatoshis(neededAmount)
	var utxos []metrix.RawTxInputs
	var minUTXOsSum decimal.Decimal
	for _, utxo := range *metrixresp {
		minUTXOsSum = minUTXOsSum.Add(utxo.Satoshis)
		utxos = append(utxos, metrix.RawTxInputs{TxID: utxo.TXID, Vout: utxo.OutputIndex})
		if minUTXOsSum.GreaterThanOrEqual(minimumSum) {
			return utxos, minUTXOsSum, nil
		}
	}

	return nil, decimal.Decimal{}, fmt.Errorf("Insufficient UTXO value attempted to be sent")
}

func calculateChange(balance, neededAmount decimal.Decimal) (decimal.Decimal, error) {
	if balance.LessThan(neededAmount) {
		return decimal.Decimal{}, fmt.Errorf("insufficient funds to create fee to chain")
	}
	return balance.Sub(neededAmount), nil
}

func calculateNeededAmount(value, gasLimit, gasPrice decimal.Decimal) decimal.Decimal {
	return value.Add(gasLimit.Mul(gasPrice))
}

func (p *ProxyETHSignTransaction) requestSendToContract(ethtx *eth.SendTransactionRequest) (string, error) {
	gasLimit, gasPrice, err := EthGasToMetrix(ethtx)
	if err != nil {
		return "", err
	}

	amount := decimal.NewFromFloat(0.0)
	if ethtx.Value != "" {
		var err error
		amount, err = EthValueToMetrixAmount(ethtx.Value, ZeroSatoshi)
		if err != nil {
			return "", errors.Wrap(err, "EthValueToMetrixAmount:")
		}
	}

	newGasPrice, err := decimal.NewFromString(gasPrice)
	if err != nil {
		return "", err
	}
	neededAmount := calculateNeededAmount(amount, decimal.NewFromBigInt(gasLimit, 0), newGasPrice)

	inputs, balance, err := p.getRequiredUtxos(ethtx.From, neededAmount)
	if err != nil {
		return "", err
	}

	change, err := calculateChange(balance, neededAmount)
	if err != nil {
		return "", err
	}

	contractInteractTx := &metrix.SendToContractRawRequest{
		ContractAddress: utils.RemoveHexPrefix(ethtx.To),
		Datahex:         utils.RemoveHexPrefix(ethtx.Data),
		Amount:          amount,
		GasLimit:        gasLimit,
		GasPrice:        gasPrice,
	}

	if from := ethtx.From; from != "" && utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return "", err
		}
		contractInteractTx.SenderAddress = from
	}

	fromAddr := utils.RemoveHexPrefix(ethtx.From)

	acc := p.Metrix.Accounts.FindByHexAddress(strings.ToLower(fromAddr))
	if acc == nil {
		return "", errors.Errorf("No such account: %s", fromAddr)
	}

	rawtxreq := []interface{}{inputs, []interface{}{map[string]*metrix.SendToContractRawRequest{"contract": contractInteractTx}, map[string]decimal.Decimal{contractInteractTx.SenderAddress: change}}}
	var rawTx string
	if err := p.Metrix.Request(metrix.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *metrix.SignRawTxResponse
	if err := p.Metrix.Request(metrix.MethodSignRawTx, []interface{}{rawTx}, &resp); err != nil {
		return "", err
	}
	if !resp.Complete {
		return "", fmt.Errorf("something went wrong with signing the transaction; transaction incomplete")
	}
	return utils.AddHexPrefix(resp.Hex), nil
}

func (p *ProxyETHSignTransaction) requestSendToAddress(req *eth.SendTransactionRequest) (string, error) {
	getMetrixWalletAddress := func(addr string) (string, error) {
		if utils.IsEthHexAddress(addr) {
			return p.FromHexAddress(utils.RemoveHexPrefix(addr))
		}
		return addr, nil
	}

	to, err := getMetrixWalletAddress(req.To)
	if err != nil {
		return "", err
	}

	from, err := getMetrixWalletAddress(req.From)
	if err != nil {
		return "", err
	}

	amount, err := EthValueToMetrixAmount(req.Value, ZeroSatoshi)
	if err != nil {
		return "", errors.Wrap(err, "EthValueToMetrixAmount:")
	}

	inputs, balance, err := p.getRequiredUtxos(req.From, amount)
	if err != nil {
		return "", err
	}

	change, err := calculateChange(balance, amount)
	if err != nil {
		return "", err
	}

	var addressValMap = map[string]decimal.Decimal{to: amount, from: change}
	rawtxreq := []interface{}{inputs, addressValMap}
	var rawTx string
	if err := p.Metrix.Request(metrix.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *metrix.SignRawTxResponse
	signrawtxreq := []interface{}{rawTx}
	if err := p.Metrix.Request(metrix.MethodSignRawTx, signrawtxreq, &resp); err != nil {
		return "", err
	}
	if !resp.Complete {
		return "", fmt.Errorf("something went wrong with signing the transaction; transaction incomplete")
	}
	return utils.AddHexPrefix(resp.Hex), nil
}

func (p *ProxyETHSignTransaction) requestCreateContract(req *eth.SendTransactionRequest) (string, error) {
	gasLimit, gasPrice, err := EthGasToMetrix(req)
	if err != nil {
		return "", err
	}

	from := req.From
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return "", err
		}
	}

	contractDeploymentTx := &metrix.CreateContractRawRequest{
		ByteCode:      utils.RemoveHexPrefix(req.Data),
		GasLimit:      gasLimit,
		GasPrice:      gasPrice,
		SenderAddress: from,
	}

	newGasPrice, err := decimal.NewFromString(gasPrice)
	if err != nil {
		return "", err
	}
	neededAmount := calculateNeededAmount(decimal.NewFromFloat(0.0), decimal.NewFromBigInt(gasLimit, 0), newGasPrice)

	inputs, balance, err := p.getRequiredUtxos(req.From, neededAmount)
	if err != nil {
		return "", err
	}

	change, err := calculateChange(balance, neededAmount)
	if err != nil {
		return "", err
	}

	rawtxreq := []interface{}{inputs, []interface{}{map[string]*metrix.CreateContractRawRequest{"contract": contractDeploymentTx}, map[string]decimal.Decimal{from: change}}}
	var rawTx string
	if err := p.Metrix.Request(metrix.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *metrix.SignRawTxResponse
	signrawtxreq := []interface{}{rawTx}
	if err := p.Metrix.Request(metrix.MethodSignRawTx, signrawtxreq, &resp); err != nil {
		return "", err
	}
	if !resp.Complete {
		return "", fmt.Errorf("something went wrong with signing the transaction; transaction incomplete")
	}
	return utils.AddHexPrefix(resp.Hex), nil
}
