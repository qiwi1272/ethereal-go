package ethereal

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var DomainHash []byte // precomputed by InitDomain

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

func Scale1e9(s string) (*big.Int, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return nil, fmt.Errorf("bad decimal %q", s)
	}
	r.Mul(r, big.NewRat(1_000_000_000, 1))
	n := new(big.Int)
	n.Div(r.Num(), r.Denom())
	return n, nil
}

type SignedMessage[T Signable] struct {
	Data      T      `json:"data"`
	Signature string `json:"signature"`
}

type Signer interface {
	getPk() *ecdsa.PrivateKey
	GetTypes() *abi.TypedData
}

type Signable interface {
	build(cl SubaccountHolder)
	ToMessage() (abi.TypedDataMessage, error)
}

func Sign(message Signable, primaryType string, signer Signer) (string, error) {
	msg, err := message.ToMessage()
	if err != nil {
		return "", err
	}

	messageHash, err := signer.GetTypes().HashStruct(primaryType, msg)
	if err != nil {
		return "", err
	}

	fullHash := MakeFullHash(DomainHash, messageHash)

	sig, err := crypto.Sign(fullHash, signer.getPk())
	if err != nil {
		return "", err
	}
	sig[64] += 27 // recovery byte fix  |  much love to _0xmer :)

	return "0x" + hex.EncodeToString(sig), nil
}

func MakeFullHash(domainHash []byte, messageHash []byte) []byte {
	fullHash := make([]byte, 0, 66)
	fullHash = append(fullHash, 0x19, 0x01)
	fullHash = append(fullHash, domainHash...)
	fullHash = append(fullHash, messageHash...)
	return crypto.Keccak256(fullHash)
}
