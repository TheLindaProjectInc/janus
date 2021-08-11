package transformer

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

// ProxyETHGetTransactionByBlockNumberAndIndex implements ETHProxy
type ProxyETHGetTransactionByBlockNumberAndIndex struct {
	*metrix.Metrix
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) Method() string {
	return "eth_getTransactionByBlockNumberAndIndex"
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.GetTransactionByBlockNumberAndIndex
	if err := json.Unmarshal(rawreq.Params, &req); err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal request")
	}
	if req.BlockNumber == "" {
		return nil, errors.New("invalid argument 0: empty hex string")
	}

	return p.request(&req)
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) request(req *eth.GetTransactionByBlockNumberAndIndex) (interface{}, error) {
	// Decoded by ProxyETHGetTransactionByBlockHashAndIndex, quickly decode so we can fail cheaply without making any calls
	_, err := hexutil.DecodeUint64(req.TransactionIndex)
	if err != nil {
		return nil, errors.Wrap(err, "invalid argument 1")
	}

	blockNum, err := getBlockNumberByParam(p.Metrix, req.BlockNumber, false)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get block number by parameter")
	}

	blockHash, err := proxyETHGetBlockByHash(p, p.Metrix, blockNum)
	if err != nil {
		return nil, err
	}
	if blockHash == nil {
		return nil, nil
	}

	var (
		getBlockByHashReq = &eth.GetTransactionByBlockHashAndIndex{
			BlockHash:        string(*blockHash),
			TransactionIndex: req.TransactionIndex,
		}
		proxy = &ProxyETHGetTransactionByBlockHashAndIndex{Metrix: p.Metrix}
	)
	return proxy.request(getBlockByHashReq)
}
