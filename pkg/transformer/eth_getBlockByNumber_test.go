package transformer

import (
	"encoding/json"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/internal"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

func initializeProxyETHGetBlockByNumber(metrixClient *metrix.Metrix) ETHProxy {
	return &ProxyETHGetBlockByNumber{metrixClient}
}

func TestGetBlockByNumberRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`false`)},
		&internal.GetTransactionByHashResponse,
	)
}

func TestGetBlockByNumberWithTransactionsRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`true`)},
		&internal.GetTransactionByHashResponseWithTransactions,
	)
}

func TestGetBlockByNumberUnknownBlockRequest(t *testing.T) {
	requestParams := []json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`true`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	metrixClient, err := internal.CreateMockedClient(mockedClientDoer)

	unknownBlockResponse := metrix.GetErrorResponse(metrix.ErrInvalidParameter)
	err = mockedClientDoer.AddError(metrix.MethodGetBlockHash, unknownBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBlockByNumber{metrixClient}
	got, err := proxyEth.Request(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	if got != (*eth.GetBlockByNumberResponse)(nil) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			string("nil"),
			string(internal.MustMarshalIndent(got, "", "  ")),
		)
	}
}
