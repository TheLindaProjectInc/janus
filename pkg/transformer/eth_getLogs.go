package transformer

import (
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/conversion"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

// ProxyETHGetLogs implements ETHProxy
type ProxyETHGetLogs struct {
	*metrix.Metrix
}

func (p *ProxyETHGetLogs) Method() string {
	return "eth_getLogs"
}

func (p *ProxyETHGetLogs) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.GetLogsRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	// TODO: Graph Node is sending the topic
	// if len(req.Topics) != 0 {
	// 	return nil, errors.New("topics is not supported yet")
	// }

	// Calls ToRequest in order transform ETH-Request to a Metrix-Request
	metrixreq, err := p.ToRequest(&req)
	if err != nil {
		return nil, err
	}

	return p.request(metrixreq)
}

func (p *ProxyETHGetLogs) request(req *metrix.SearchLogsRequest) (*eth.GetLogsResponse, error) {
	receipts, err := p.SearchLogs(req)
	if err != nil {
		return nil, err
	}

	logs := make([]eth.Log, 0)
	for _, receipt := range receipts {
		r := metrix.TransactionReceipt(receipt)
		logs = append(logs, conversion.ExtractETHLogsFromTransactionReceipt(&r)...)
	}

	resp := eth.GetLogsResponse(logs)
	return &resp, nil
}

func (p *ProxyETHGetLogs) ToRequest(ethreq *eth.GetLogsRequest) (*metrix.SearchLogsRequest, error) {
	//transform EthRequest fromBlock to MetrixReq fromBlock:
	from, err := getBlockNumberByRawParam(p.Metrix, ethreq.FromBlock, true)
	if err != nil {
		return nil, err
	}

	//transform EthRequest toBlock to MetrixReq toBlock:
	to, err := getBlockNumberByRawParam(p.Metrix, ethreq.ToBlock, true)
	if err != nil {
		return nil, err
	}

	//transform EthReq address to MetrixReq address:
	var addresses []string
	if ethreq.Address != nil {
		if isBytesOfString(ethreq.Address) {
			var addr string
			if err = json.Unmarshal(ethreq.Address, &addr); err != nil {
				return nil, err
			}
			addresses = append(addresses, addr)
		} else {
			if err = json.Unmarshal(ethreq.Address, &addresses); err != nil {
				return nil, err
			}
		}
		for i := range addresses {
			addresses[i] = utils.RemoveHexPrefix(addresses[i])
		}
	}

	//transform EthReq topics to MetrixReq topics:
	topics, err := eth.TranslateTopics(ethreq.Topics)
	if err != nil {
		return nil, err
	}

	return &metrix.SearchLogsRequest{
		Addresses: addresses,
		FromBlock: from,
		ToBlock:   to,
		Topics:    topics,
	}, nil
}
