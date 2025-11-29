package ethereal

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var domainHash []byte // precomputed by InitDomain

// ------- HELPERS -------
func ParseTypeSchema(typeString string) ([]abi.Type, error) {
	fields := strings.Split(typeString, ",")
	args := make([]abi.Type, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		// split into "<type> <name>"
		parts := strings.Fields(field)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field %q in type string", field)
		}
		arg := abi.Type{
			Name: parts[1],
			Type: parts[0],
		}
		args = append(args, arg)
	}
	return args, nil
}

func scale1e9(s string) (*big.Int, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return nil, fmt.Errorf("bad decimal %q", s)
	}
	r.Mul(r, big.NewRat(1_000_000_000, 1))
	n := new(big.Int)
	n.Div(r.Num(), r.Denom())
	return n, nil
}

func strip0x(s string) string {
	if len(s) > 1 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

func (lo *LimitOrder) toMessage() (abi.TypedDataMessage, error) {
	qtyBig, err := scale1e9(lo.Quantity)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}
	priceBig, err := scale1e9(lo.Price)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}

	// even though we expect these values to be uint8 according to their signatures,
	// setting them as native uint8 raises a compiler error. strings or big ints are accepted.
	side := new(big.Int).SetInt64(lo.Side)
	engine := new(big.Int).SetInt64(lo.EngineType)
	id := new(big.Int).SetInt64(lo.OnchainID)
	sigTs := new(big.Int).SetInt64(lo.SignedAt)

	return abi.TypedDataMessage{
		"sender":     lo.Sender,
		"subaccount": lo.Subaccount,
		"quantity":   qtyBig,
		"price":      priceBig,
		"reduceOnly": lo.ReduceOnly,
		"side":       side,
		"engineType": engine,
		"productId":  id,
		"nonce":      lo.Nonce,
		"signedAt":   sigTs,
	}, nil
}

func (co *CancelOrder) toMessage() (abi.TypedDataMessage, error) {
	return abi.TypedDataMessage{
		"sender":     co.Sender,
		"subaccount": co.Subaccount,
		"nonce":      co.Nonce,
		//"orderIds":   co.OrderIDs,
	}, nil
}

type Signable interface {
	*LimitOrder | *CancelOrder
	toMessage() (abi.TypedDataMessage, error)
}

func Sign[T Signable](message T, primaryType string, cl *EtherealClient) (string, error) {
	msg, err := message.toMessage()
	if err != nil {
		return "", err
	}

	messageHash, err := cl.types.HashStruct(primaryType, msg)
	if err != nil {
		return "", err
	}

	fullHash := make([]byte, 0, 66)
	fullHash = append(fullHash, 0x19, 0x01)
	fullHash = append(fullHash, domainHash...)
	fullHash = append(fullHash, messageHash...)

	fullHash = crypto.Keccak256(fullHash)

	sig, err := crypto.Sign(fullHash, cl.pk)
	if err != nil {
		return "", err
	}
	sig[64] += 27 // recovery byte fix  |  much love to _0xmer :)

	return "0x" + hex.EncodeToString(sig), nil
}
