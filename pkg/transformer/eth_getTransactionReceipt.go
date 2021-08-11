package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/conversion"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

var STATUS_SUCCESS = "0x1"
var STATUS_FAILURE = "0x0"

// ProxyETHGetTransactionReceipt implements ETHProxy
type ProxyETHGetTransactionReceipt struct {
	*metrix.Metrix
}

func (p *ProxyETHGetTransactionReceipt) Method() string {
	return "eth_getTransactionReceipt"
}

func (p *ProxyETHGetTransactionReceipt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.GetTransactionReceiptRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}
	if req == "" {
		return nil, errors.New("empty transaction hash")
	}
	var (
		txHash  = utils.RemoveHexPrefix(string(req))
		metrixReq = metrix.GetTransactionReceiptRequest(txHash)
	)
	return p.request(&metrixReq)
}

func (p *ProxyETHGetTransactionReceipt) request(req *metrix.GetTransactionReceiptRequest) (*eth.GetTransactionReceiptResponse, error) {
	metrixReceipt, err := p.Metrix.GetTransactionReceipt(string(*req))
	if err != nil {
		ethTx, getRewardTransactionErr := getRewardTransactionByHash(p.Metrix, string(*req))
		if getRewardTransactionErr != nil {
			errCause := errors.Cause(err)
			if errCause == metrix.EmptyResponseErr {
				return nil, nil
			}
			p.Metrix.GetDebugLogger().Log("msg", "Transaction does not exist", "txid", string(*req))
			return nil, err
		}
		return &eth.GetTransactionReceiptResponse{
			TransactionHash:   ethTx.Hash,
			TransactionIndex:  ethTx.TransactionIndex,
			BlockHash:         ethTx.BlockHash,
			BlockNumber:       ethTx.BlockNumber,
			CumulativeGasUsed: "0x0",
			GasUsed:           "0x0",
			From:              ethTx.From,
			To:                ethTx.To,
			Logs:              []eth.Log{},
			LogsBloom:         eth.EmptyLogsBloom,
			Status:            STATUS_SUCCESS,
		}, nil
	}

	ethReceipt := &eth.GetTransactionReceiptResponse{
		TransactionHash:   utils.AddHexPrefix(metrixReceipt.TransactionHash),
		TransactionIndex:  hexutil.EncodeUint64(metrixReceipt.TransactionIndex),
		BlockHash:         utils.AddHexPrefix(metrixReceipt.BlockHash),
		BlockNumber:       hexutil.EncodeUint64(metrixReceipt.BlockNumber),
		ContractAddress:   utils.AddHexPrefixIfNotEmpty(metrixReceipt.ContractAddress),
		CumulativeGasUsed: hexutil.EncodeUint64(metrixReceipt.CumulativeGasUsed),
		GasUsed:           hexutil.EncodeUint64(metrixReceipt.GasUsed),
		From:              utils.AddHexPrefixIfNotEmpty(metrixReceipt.From),
		To:                utils.AddHexPrefixIfNotEmpty(metrixReceipt.To),

		// TODO: researching
		// ! Temporary accept this value to be always zero, as it is at eth logs
		LogsBloom: eth.EmptyLogsBloom,
	}

	status := STATUS_FAILURE
	if metrixReceipt.Excepted == "None" {
		status = STATUS_SUCCESS
	} else {
		p.Metrix.GetDebugLogger().Log("transaction", ethReceipt.TransactionHash, "msg", "transaction excepted", "message", metrixReceipt.Excepted)
	}
	ethReceipt.Status = status

	r := metrix.TransactionReceipt(*metrixReceipt)
	ethReceipt.Logs = conversion.ExtractETHLogsFromTransactionReceipt(&r)

	metrixTx, err := p.Metrix.GetRawTransaction(metrixReceipt.TransactionHash, false)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get transaction")
	}
	decodedRawMetrixTx, err := p.Metrix.DecodeRawTransaction(metrixTx.Hex)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't decode raw transaction")
	}
	if decodedRawMetrixTx.IsContractCreation() {
		ethReceipt.To = ""
	} else {
		ethReceipt.ContractAddress = ""
	}

	// TODO: researching
	// - The following code reason is unknown (see original comment)
	// - Code temporary commented, until an error occures
	// ! Do not remove
	// // contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
	// if status != "0x1" {
	// 	// if failure, should return null for contractAddress, instead of the zero address.
	// 	ethTxReceipt.ContractAddress = ""
	// }

	return ethReceipt, nil
}
