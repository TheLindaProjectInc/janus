package transformer

import (
	"encoding/json"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/internal"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

func initializeProxyETHGetTransactionByBlockHashAndIndex(metrixClient *metrix.Metrix) ETHProxy {
	return &ProxyETHGetTransactionByBlockHashAndIndex{metrixClient}
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockHashAndIndex,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockHash + `"`), []byte(`"0x0"`)},
		internal.GetTransactionByHashResponseData,
	)
}
