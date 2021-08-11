package transformer

import (
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/notifier"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"
)

type Transformer struct {
	metrixClient   *metrix.Metrix
	debugMode    bool
	logger       log.Logger
	transformers map[string]ETHProxy
}

// New creates a new Transformer
func New(metrixClient *metrix.Metrix, proxies []ETHProxy, opts ...Option) (*Transformer, error) {
	if metrixClient == nil {
		return nil, errors.New("metrixClient cannot be nil")
	}

	t := &Transformer{
		metrixClient: metrixClient,
		logger:     log.NewNopLogger(),
	}

	var err error
	for _, p := range proxies {
		if err = t.Register(p); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Register registers an ETHProxy to a Transformer
func (t *Transformer) Register(p ETHProxy) error {
	if t.transformers == nil {
		t.transformers = make(map[string]ETHProxy)
	}

	m := p.Method()
	if _, ok := t.transformers[m]; ok {
		return errors.Errorf("method already exist: %s ", m)
	}

	t.transformers[m] = p

	return nil
}

// Transform takes a Transformer and transforms the request from ETH request and returns the proxy request
func (t *Transformer) Transform(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	proxy, err := t.getProxy(req.Method)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get proxy")
	}
	resp, err := proxy.Request(req, c)
	if err != nil {
		return nil, errors.WithMessagef(err, "couldn't proxy %s request", req.Method)
	}
	return resp, nil
}

func (t *Transformer) getProxy(method string) (ETHProxy, error) {
	proxy, ok := t.transformers[method]
	if !ok {
		return nil, errors.Errorf("The method %s does not exist/is not available", method)
	}
	return proxy, nil
}

func (t *Transformer) IsDebugEnabled() bool {
	return t.debugMode
}

// DefaultProxies are the default proxy methods made available
func DefaultProxies(metrixRPCClient *metrix.Metrix, agent *notifier.Agent) []ETHProxy {
	filter := eth.NewFilterSimulator()
	getFilterChanges := &ProxyETHGetFilterChanges{Metrix: metrixRPCClient, filter: filter}
	ethCall := &ProxyETHCall{Metrix: metrixRPCClient}

	return []ETHProxy{
		ethCall,
		&ProxyNetListening{Metrix: metrixRPCClient},
		&ProxyETHPersonalUnlockAccount{},
		&ProxyETHChainId{Metrix: metrixRPCClient},
		&ProxyETHBlockNumber{Metrix: metrixRPCClient},
		&ProxyETHHashrate{Metrix: metrixRPCClient},
		&ProxyETHMining{Metrix: metrixRPCClient},
		&ProxyETHNetVersion{Metrix: metrixRPCClient},
		&ProxyETHGetTransactionByHash{Metrix: metrixRPCClient},
		&ProxyETHGetTransactionByBlockNumberAndIndex{Metrix: metrixRPCClient},
		&ProxyETHGetLogs{Metrix: metrixRPCClient},
		&ProxyETHGetTransactionReceipt{Metrix: metrixRPCClient},
		&ProxyETHSendTransaction{Metrix: metrixRPCClient},
		&ProxyETHAccounts{Metrix: metrixRPCClient},
		&ProxyETHGetCode{Metrix: metrixRPCClient},

		&ProxyETHNewFilter{Metrix: metrixRPCClient, filter: filter},
		&ProxyETHNewBlockFilter{Metrix: metrixRPCClient, filter: filter},
		getFilterChanges,
		&ProxyETHGetFilterLogs{ProxyETHGetFilterChanges: getFilterChanges},
		&ProxyETHUninstallFilter{Metrix: metrixRPCClient, filter: filter},

		&ProxyETHEstimateGas{ProxyETHCall: ethCall},
		&ProxyETHGetBlockByNumber{Metrix: metrixRPCClient},
		&ProxyETHGetBlockByHash{Metrix: metrixRPCClient},
		&ProxyETHGetBalance{Metrix: metrixRPCClient},
		&ProxyETHGetStorageAt{Metrix: metrixRPCClient},
		&ETHGetCompilers{},
		&ETHProtocolVersion{},
		&ETHGetUncleByBlockHashAndIndex{},
		&ETHGetUncleCountByBlockHash{},
		&ETHGetUncleCountByBlockNumber{},
		&Web3ClientVersion{},
		&Web3Sha3{},
		&ProxyETHSign{Metrix: metrixRPCClient},
		&ProxyETHGasPrice{Metrix: metrixRPCClient},
		&ProxyETHTxCount{Metrix: metrixRPCClient},
		&ProxyETHSignTransaction{Metrix: metrixRPCClient},
		&ProxyETHSendRawTransaction{Metrix: metrixRPCClient},

		&ETHSubscribe{Metrix: metrixRPCClient, Agent: agent},
		&ETHUnsubscribe{Metrix: metrixRPCClient, Agent: agent},

		&ProxyMETRIXGetUTXOs{Metrix: metrixRPCClient},

		&ProxyNetPeerCount{Metrix: metrixRPCClient},
	}
}

func SetDebug(debug bool) func(*Transformer) error {
	return func(t *Transformer) error {
		t.debugMode = debug
		return nil
	}
}

func SetLogger(l log.Logger) func(*Transformer) error {
	return func(t *Transformer) error {
		t.logger = log.WithPrefix(l, "component", "transformer")
		return nil
	}
}
