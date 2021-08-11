package metrix

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type Accounts []*btcutil.WIF

func (as Accounts) FindByHexAddress(addr string) *btcutil.WIF {
	for _, a := range as {
		acc := &Account{a}

		if addr == acc.ToHexAddress() {
			return a
		}
	}

	return nil
}

type Account struct {
	*btcutil.WIF
}

func (a *Account) ToHexAddress() string {
	// wif := (*btcutil.WIF)(a)

	keyid := btcutil.Hash160(a.SerializePubKey())
	return hex.EncodeToString(keyid)
}

var metrixMainNetParams = chaincfg.MainNetParams
var metrixTestNetParams = chaincfg.MainNetParams

func init() {
	metrixMainNetParams.PubKeyHashAddrID = 50
	metrixMainNetParams.ScriptHashAddrID = 85

	metrixTestNetParams.PubKeyHashAddrID = 110
	metrixTestNetParams.ScriptHashAddrID = 187
}

func (a *Account) ToBase58Address(isMain bool) (string, error) {
	params := &metrixMainNetParams
	if !isMain {
		params = &metrixTestNetParams
	}

	addr, err := btcutil.NewAddressPubKey(a.SerializePubKey(), params)
	if err != nil {
		return "", err
	}

	return addr.AddressPubKeyHash().String(), nil
}
