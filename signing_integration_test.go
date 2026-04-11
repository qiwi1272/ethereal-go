// Optional live Testnet checks: set ETHEREAL_INTEGRATION=1 (and ETHEREAL_PK if required).
package etherealRest

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func TestLive_orderSigning_againstTestnet(t *testing.T) {
	if os.Getenv("ETHEREAL_INTEGRATION") != "1" {
		t.Skip("set ETHEREAL_INTEGRATION=1 to run live Testnet signing checks")
	}
	pk := os.Getenv("ETHEREAL_PK")
	if pk == "" {
		pk = "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a"
	}

	order := Order{
		Sender:     "0xdeadbeef00000000000000000000000000000000",
		Subaccount: "0x123456789abcde00000000000000000000000000000000000000000000000000",
		Quantity:   "1",
		Price:      "3000",
		ReduceOnly: false,
		Side:       BUY,
		EngineType: PERPETUAL,
		OnchainID:  2,
		Nonce:      "1764897077655477722",
		SignedAt:   int64(1764897077),
	}

	ctx := context.Background()
	client, err := NewClient(ctx, pk, Testnet)
	if err != nil {
		t.Fatal(err)
	}

	domainHashString, err := client.InitDomain(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expectedDomainHash := "baf501bc2614cf7092d082742580b04c176be1815f46e407eab1bc37ba543c05"
	if domainHashString != expectedDomainHash {
		t.Fatalf("domain hash drift (API may have changed EIP-712 domain): got %s want %s", domainHashString, expectedDomainHash)
	}

	msg, err := order.ToMessage()
	if err != nil {
		t.Fatal(err)
	}
	messageHash, err := client.GetTypes().HashStruct("TradeOrder", msg)
	if err != nil {
		t.Fatal(err)
	}

	signer := &integrationTestSigner{pk: mustIntegrationKey(t, pk), types: client.GetTypes()}

	signature, err := Sign(&order, "TradeOrder", signer)
	if err != nil {
		t.Fatal(err)
	}

	domainBytes, err := reverseHexIntegration(domainHashString)
	if err != nil {
		t.Fatal(err)
	}
	fullHash := MakeFullHash(domainBytes, messageHash)
	t.Logf("message hash: %s", common.Bytes2Hex(messageHash))
	t.Logf("full hash: %s", common.Bytes2Hex(fullHash))

	expectedSignature := "0x82aed7486e9855459f58537e413760597e689d3ba7b859f56b6edc730e044fff2888ccf92cd282a8299d8d6a76f8bf0aa93d97f30340c4bb0d27b626aca62f211b"
	if signature != expectedSignature {
		t.Fatalf("signature drift: got %s want %s", signature, expectedSignature)
	}

	payload := SignedMessage[*Order]{Data: &order, Signature: signature}
	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signed payload:\n%s", string(payloadJSON))
}

type integrationTestSigner struct {
	pk    *ecdsa.PrivateKey
	types *abi.TypedData
}

func mustIntegrationKey(t *testing.T, pk string) *ecdsa.PrivateKey {
	t.Helper()
	k, err := crypto.HexToECDSA(strings.TrimPrefix(pk, "0x"))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func (s *integrationTestSigner) GetPk() *ecdsa.PrivateKey {
	return s.pk
}

func (s *integrationTestSigner) GetTypes() *abi.TypedData {
	return s.types
}

func reverseHexIntegration(s string) ([]byte, error) {
	clean := strings.TrimPrefix(s, "0x")
	return hex.DecodeString(clean)
}
