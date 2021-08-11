package transformer

import (
	"github.com/labstack/echo"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
)

// Web3ClientVersion implements web3_clientVersion
type Web3ClientVersion struct {
	// *metrix.Metrix
}

func (p *Web3ClientVersion) Method() string {
	return "web3_clientVersion"
}

func (p *Web3ClientVersion) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return "METRIX ETHTestRPC/ethereum-js", nil
}

// func (p *Web3ClientVersion) ToResponse(ethresp *metrix.CallContractResponse) *eth.CallResponse {
// 	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
// 	metrixresp := eth.CallResponse(data)
// 	return &metrixresp
// }
