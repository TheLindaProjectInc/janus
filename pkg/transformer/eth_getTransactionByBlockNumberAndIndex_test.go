package transformer

import (
	"encoding/json"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/internal"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

func initializeProxyETHGetTransactionByBlockNumberAndIndex(metrixClient *metrix.Metrix) ETHProxy {
	return &ProxyETHGetTransactionByBlockNumberAndIndex{metrixClient}
}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockNumberAndIndex,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x0"`)},
		internal.GetTransactionByHashResponseData,
	)
}
