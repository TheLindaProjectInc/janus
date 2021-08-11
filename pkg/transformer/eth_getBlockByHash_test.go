package transformer

import (
	"encoding/json"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/internal"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

func initializeProxyETHGetBlockByHash(metrixClient *metrix.Metrix) ETHProxy {
	return &ProxyETHGetBlockByHash{metrixClient}
}

func TestGetBlockByHashRequestNonceLength(t *testing.T) {
	if len(utils.RemoveHexPrefix(internal.GetTransactionByHashResponse.Nonce)) != 16 {
		t.Errorf("Nonce test data should be zero left padded length 16")
	}
}

func TestGetBlockByHashRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockHexHash + `"`), []byte(`false`)},
		&internal.GetTransactionByHashResponse,
	)
}

func TestGetBlockByHashTransactionsRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockHexHash + `"`), []byte(`true`)},
		&internal.GetTransactionByHashResponseWithTransactions,
	)
}
