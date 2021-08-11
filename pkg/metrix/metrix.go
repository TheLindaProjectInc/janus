package metrix

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/TheLindaProjectInc/janus/pkg/utils"
)

type Metrix struct {
	*Client
	*Method
	chain string
}

const (
	ChainMain    = "main"
	ChainTest    = "test"
	ChainRegTest = "regtest"
)

var AllChains = []string{ChainMain, ChainRegTest, ChainTest}

func New(c *Client, chain string) (*Metrix, error) {
	if !utils.InStrSlice(AllChains, chain) {
		return nil, errors.New("invalid metrix chain")
	}

	return &Metrix{
		Client: c,
		Method: &Method{Client: c},
		chain:  chain,
	}, nil
}

func (c *Metrix) Chain() string {
	return c.chain
}

// Presents hexed address prefix of a specific chain without
// `0x` prefix, this is a ready to use hexed string
type HexAddressPrefix string

const (
	PrefixMainChainAddress    HexAddressPrefix = "32"
	PrefixTestChainAddress    HexAddressPrefix = "6E"
	PrefixRegTestChainAddress HexAddressPrefix = PrefixTestChainAddress
)

// Returns decoded hexed string prefix, as ready to use slice of bytes
func (prefix HexAddressPrefix) AsBytes() ([]byte, error) {
	bytes, err := hex.DecodeString(string(prefix))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't decode hexed string")
	}
	return bytes, nil
}

// Returns first 4 bytes of a double sha256 hash of the provided `prefixedAddrBytes`,
// which must be already prefixed with a specific chain prefix
func CalcAddressChecksum(prefixedAddr []byte) []byte {
	hash := sha256.Sum256(prefixedAddr)
	hash = sha256.Sum256(hash[:])
	return hash[:4]
}
