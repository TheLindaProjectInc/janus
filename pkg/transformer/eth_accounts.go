package transformer

import (
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

// ProxyETHAccounts implements ETHProxy
type ProxyETHAccounts struct {
	*metrix.Metrix
}

func (p *ProxyETHAccounts) Method() string {
	return "eth_accounts"
}

func (p *ProxyETHAccounts) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHAccounts) request() (eth.AccountsResponse, error) {
	var accounts eth.AccountsResponse

	for _, acc := range p.Accounts {
		acc := metrix.Account{acc}
		addr := acc.ToHexAddress()

		accounts = append(accounts, utils.AddHexPrefix(addr))
	}

	return accounts, nil
}

func (p *ProxyETHAccounts) ToResponse(ethresp *metrix.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	metrixresp := eth.CallResponse(data)
	return &metrixresp
}
