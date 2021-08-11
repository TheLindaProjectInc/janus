package transformer

import (
	"encoding/json"
	"testing"

	"github.com/TheLindaProjectInc/janus/pkg/metrix"
	"github.com/TheLindaProjectInc/janus/pkg/internal"
)

func TestNetListeningInactive(t *testing.T) {
	testNetListeningRequest(t, false)
}

func TestNetListeningActive(t *testing.T) {
	testNetListeningRequest(t, true)
}

func testNetListeningRequest(t *testing.T, active bool) {
	//preparing the request
	requestParams := []json.RawMessage{} //net_listening has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	metrixClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	networkInfoResponse := metrix.NetworkInfoResponse{NetworkActive: active}
	err = mockedClientDoer.AddResponseWithRequestID(2, metrix.MethodGetNetworkInfo, networkInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyNetListening{metrixClient}
	got, err := proxyEth.Request(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := active
	if want != got {
		t.Errorf(
			"error\nwant: %t\ngot: %t",
			want,
			got,
		)
	}
}
