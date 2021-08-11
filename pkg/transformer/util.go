package transformer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
	"github.com/TheLindaProjectInc/janus/pkg/eth"
	"github.com/TheLindaProjectInc/janus/pkg/metrix"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

var ZeroSatoshi = decimal.NewFromInt(0)
var OneSatoshi = decimal.NewFromFloat(0.00000001)
var MinimumGas = decimal.NewFromFloat(0.0000004)

type EthGas interface {
	GasHex() string
	GasPriceHex() string
}

func EthGasToMetrix(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = g.(*eth.SendTransactionRequest).Gas.Int

	gasPriceDecimal, err := EthValueToMetrixAmount(g.GasPriceHex(), MinimumGas)
	if err != nil {
		return nil, "0.0", err
	}
	if gasPriceDecimal.LessThan(MinimumGas) {
		gasPriceDecimal = MinimumGas
	}
	gasPrice = fmt.Sprintf("%v", gasPriceDecimal)

	return
}

func MetrixGasToEth(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = g.(*eth.SendTransactionRequest).Gas.Int

	gasPriceDecimal, err := EthValueToMetrixAmount(g.GasPriceHex(), MinimumGas)
	if err != nil {
		return nil, "0.0", err
	}
	if gasPriceDecimal.LessThan(MinimumGas) {
		gasPriceDecimal = MinimumGas
	}
	gasPrice = fmt.Sprintf("%v", gasPriceDecimal)

	return
}

func EthValueToMetrixAmount(val string, defaultValue decimal.Decimal) (decimal.Decimal, error) {
	if val == "" {
		return defaultValue, nil
	}

	ethVal, err := utils.DecodeBig(val)
	if err != nil {
		return ZeroSatoshi, err
	}

	ethValDecimal, err := decimal.NewFromString(ethVal.String())
	if err != nil {
		return ZeroSatoshi, errors.New("decimal.NewFromString was not a success")
	}

	return EthDecimalValueToMetrixAmount(ethValDecimal), nil
}

func EthDecimalValueToMetrixAmount(ethValDecimal decimal.Decimal) decimal.Decimal {
	// Convert Wei to Metrix
	// 10000000000
	// one satoshi is 0.00000001
	// we need to drop precision for values smaller than that
	maximumPrecision := ethValDecimal.Mul(decimal.NewFromFloat(float64(1e-8))).Floor()
	amount := maximumPrecision.Mul(decimal.NewFromFloat(float64(1e-10)))

	return amount
}

func MetrixValueToETHAmount(val string, defaultValue decimal.Decimal) (decimal.Decimal, error) {
	if val == "" {
		return defaultValue, nil
	}

	metrixVal, err := utils.DecodeBig(val)
	if err != nil {
		return ZeroSatoshi, err
	}

	metrixValDecimal, err := decimal.NewFromString(metrixVal.String())
	if err != nil {
		return ZeroSatoshi, errors.New("decimal.NewFromString was not a success")
	}

	return MetrixDecimalValueToETHAmount(metrixValDecimal), nil
}

func MetrixDecimalValueToETHAmount(metrixValDecimal decimal.Decimal) decimal.Decimal {
	// Computes inverse of EthDecimalValueToMetrixAmount
	amount := metrixValDecimal.Div(decimal.NewFromFloat(float64(1e-18)))

	return amount
}

func formatMetrixAmount(amount decimal.Decimal) (string, error) {
	decimalAmount := amount.Mul(decimal.NewFromFloat(float64(1e18)))

	//convert decimal to Integer
	result := decimalAmount.BigInt()

	if !decimalAmount.Equals(decimal.NewFromBigInt(result, 0)) {
		return "0x0", errors.New("decimal.BigInt() was not a success")
	}

	return hexutil.EncodeBig(result), nil
}

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "Invalid RPC input")
	}
	return nil
}

// NOTE:
// 	- is not for reward transactions
// 	- Vin[i].N (vout number) -> get Transaction(txID).Vout[N].Address
// 	- returning address already has 0x prefix
func getNonContractTxSenderAddress(p *metrix.Metrix, vins []*metrix.DecodedRawTransactionInV) (string, error) {
	for _, vin := range vins {
		prevMetrixTx, err := p.GetRawTransaction(vin.TxID, false)
		if err != nil {
			return "", errors.WithMessage(err, "couldn't get vin's previous transaction")
		}
		for _, out := range prevMetrixTx.Vouts {
			for _, address := range out.Details.Addresses {
				return utils.AddHexPrefix(address), nil
			}
		}
	}
	return "", errors.New("not found")
}

// NOTE:
// 	- is not for reward transactions
// 	- returning address already has 0x prefix
//
// 	TODO: researching
// 	- Vout[0].Addresses[i] != "" - temporary solution
func findNonContractTxReceiverAddress(vouts []*metrix.DecodedRawTransactionOutV) (string, error) {
	for _, vout := range vouts {
		for _, address := range vout.ScriptPubKey.Addresses {
			if address != "" {
				return utils.AddHexPrefix(address), nil
			}
		}
	}
	return "", errors.New("not found")
}

func getBlockNumberByHash(p *metrix.Metrix, hash string) (uint64, error) {
	block, err := p.GetBlock(hash)
	if err != nil {
		return 0, errors.WithMessage(err, "couldn't get block")
	}
	p.GetDebugLogger().Log("function", "getBlockNumberByHash", "hash", hash, "block", block.Height)
	return uint64(block.Height), nil
}

func getTransactionIndexInBlock(p *metrix.Metrix, txHash string, blockHash string) (int64, error) {
	block, err := p.GetBlock(blockHash)
	if err != nil {
		return -1, errors.WithMessage(err, "couldn't get block")
	}
	for i, blockTx := range block.Txs {
		if txHash == blockTx {
			p.GetDebugLogger().Log("function", "getTransactionIndexInBlock", "msg", "Found transaction index in block", "txHash", txHash, "blockHash", blockHash, "index", i)
			return int64(i), nil
		}
	}
	p.GetDebugLogger().Log("function", "getTransactionIndexInBlock", "msg", "Could not find transaction index for hash in block", "txHash", txHash, "blockHash", blockHash)
	return -1, errors.New("not found")
}

func formatMetrixNonce(nonce int) string {
	var (
		hexedNonce     = strconv.FormatInt(int64(nonce), 16)
		missedCharsNum = 16 - len(hexedNonce)
	)
	for i := 0; i < missedCharsNum; i++ {
		hexedNonce = "0" + hexedNonce
	}
	return "0x" + hexedNonce
}

// Returns Metrix block number. Result depends on a passed raw param. Raw param's slice of bytes should
// has one of the following values:
// 	- hex string representation of a number of a specific block
// 	- string "latest" - for the latest mined block
// 	- string "earliest" for the genesis block
// 	- string "pending" - for the pending state/transactions
// Uses defaultVal to differntiate from a eth_getBlockByNumber req and eth_getLogs/eth_newFilter
func getBlockNumberByRawParam(p *metrix.Metrix, rawParam json.RawMessage, defaultVal bool) (*big.Int, error) {
	if !isBytesOfString(rawParam) {
		return nil, errors.Errorf("invalid parameter format - string is expected")
	}

	param := string(rawParam[1 : len(rawParam)-1]) // trim \" runes
	return getBlockNumberByParam(p, param, defaultVal)
}

func getBlockNumberByParam(p *metrix.Metrix, param string, defaultVal bool) (*big.Int, error) {
	if len(param) < 1 {
		if defaultVal {
			res, err := p.GetBlockChainInfo()
			if err != nil {
				return nil, err
			}
			p.GetDebugLogger().Log("function", "getBlockNumberByParam", "msg", "returning default value ("+strconv.Itoa(int(res.Blocks))+")")
			return big.NewInt(res.Blocks), nil
		} else {
			return nil, errors.Errorf("empty parameter value")
		}

	}

	switch param {
	case "latest":
		res, err := p.GetBlockChainInfo()
		if err != nil {
			return nil, err
		}
		p.GetDebugLogger().Log("latest", res.Blocks, "msg", "Got latest block")
		return big.NewInt(res.Blocks), nil

	case "earliest":
		// TODO: discuss
		// ! Genesis block cannot be retreived
		return big.NewInt(0), nil

	case "pending":
		// TODO: discuss
		// 	! Researching
		return nil, errors.New("TODO: tag is in implementation")

	default: // hex number
		n, err := utils.DecodeBig(param)
		if err != nil {
			p.GetDebugLogger().Log("function", "getBlockNumberByParam", "msg", "Failed to decode hex parameter", "value", param)
			return nil, errors.Wrap(err, "couldn't decode hex number to big int")
		}
		return n, nil
	}
}

func isBytesOfString(v json.RawMessage) bool {
	dQuote := []byte{'"'}
	if !bytes.HasPrefix(v, dQuote) && !bytes.HasSuffix(v, dQuote) {
		return false
	}
	if bytes.Count(v, dQuote) != 2 {
		return false
	}
	// TODO: decide
	// ? Should we iterate over v to check if v[1:len(v)-2] is in a range of a-A, z-Z, 0-9
	return true
}

// Converts Ethereum address to a Metrix address, where `address` represents
// Ethereum address without `0x` prefix and `chain` represents target Metrix
// chain
func convertETHAddress(address string, chain string) (metrixAddress string, _ error) {
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't decode hexed address - %q", address)
	}

	var prefix []byte
	switch chain {
	case metrix.ChainMain:
		chainPrefix, err := metrix.PrefixMainChainAddress.AsBytes()
		if err != nil {
			return "", errors.WithMessagef(err, "couldn't convert %q Metrix chain prefix to slice of bytes", chain)
		}
		prefix = chainPrefix

	case metrix.ChainTest, metrix.ChainRegTest:
		chainPrefix, err := metrix.PrefixTestChainAddress.AsBytes()
		if err != nil {
			return "", errors.WithMessagef(err, "couldn't convert %q Metrix chain prefix to slice of bytes", chain)
		}
		prefix = chainPrefix

	default:
		return "", errors.Errorf("unsupported %q Metrix chain", chain)
	}

	var (
		prefixedAddrBytes = append(prefix, addrBytes...)
		checksum          = metrix.CalcAddressChecksum(prefixedAddrBytes)
		metrixAddressBytes  = append(prefixedAddrBytes, checksum...)
	)
	return base58.Encode(metrixAddressBytes), nil
}

// Converts Metrix address to an Ethereum address
func convertMetrixAddress(address string) (ethAddress string, _ error) {
	if n := len(address); n < 22 {
		return "", errors.Errorf("invalid address: length is less than 22 bytes - %d", n)
	}

	// Drop Metrix chain prefix and checksum suffix
	ethAddrBytes := base58.Decode(address)[1:21]

	return hex.EncodeToString(ethAddrBytes), nil
}

func processFilter(p *ProxyETHGetFilterChanges, rawreq *eth.JSONRPCRequest) (*eth.Filter, error) {
	var req eth.GetFilterChangesRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	filterID, err := hexutil.DecodeUint64(string(req))
	if err != nil {
		return nil, err
	}

	_filter, ok := p.filter.Filter(filterID)
	if !ok {
		return nil, errors.New("Invalid filter id")
	}
	filter := _filter.(*eth.Filter)

	return filter, nil
}

// Converts a satoshis to metrix balance
func convertFromSatoshisToMetrix(inSatoshis decimal.Decimal) decimal.Decimal {
	return inSatoshis.Div(decimal.NewFromFloat(float64(1e8)))
}

// Converts a metrix balance to satoshis
func convertFromMetrixToSatoshis(inMetrix decimal.Decimal) decimal.Decimal {
	return inMetrix.Mul(decimal.NewFromFloat(float64(1e8)))
}
