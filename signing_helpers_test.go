package etherealRest

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func TestScale1e9(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"1", "1000000000"},
		{"0.000000001", "1"},
		{"3000", "3000000000000"},
		{"1.5", "1500000000"},
	}
	for _, tc := range cases {
		got, err := Scale1e9(tc.in)
		if err != nil {
			t.Fatalf("Scale1e9(%q): %v", tc.in, err)
		}
		want := new(big.Int)
		if _, ok := want.SetString(tc.want, 10); !ok {
			t.Fatal("bad want")
		}
		if got.Cmp(want) != 0 {
			t.Fatalf("Scale1e9(%q) = %s want %s", tc.in, got.String(), want.String())
		}
	}
	if _, err := Scale1e9("not-a-number"); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseTypeSchema(t *testing.T) {
	s := "address sender, bytes32 subaccount, uint128 quantity"
	got, err := ParseTypeSchema(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Name != "sender" || got[0].Type != "address" {
		t.Fatalf("got %#v", got)
	}
	if _, err := ParseTypeSchema("invalid"); err == nil {
		t.Fatal("expected error")
	}
}

func TestMakeFullHash(t *testing.T) {
	domain, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	msg, _ := hex.DecodeString("aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899")
	h := MakeFullHash(domain, msg)
	want := crypto.Keccak256(append(append([]byte{0x19, 0x01}, domain...), msg...))
	if !strings.EqualFold(hex.EncodeToString(h), hex.EncodeToString(want)) {
		t.Fatalf("hash mismatch")
	}
}

func TestSignedMessage_JSON_roundTrip(t *testing.T) {
	o := &Order{Quantity: "1", Price: "1"}
	payload := SignedMessage[*Order]{Data: o, Signature: "0xabc"}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var out SignedMessage[*Order]
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Signature != "0xabc" || out.Data.Quantity != "1" {
		t.Fatalf("%+v", out)
	}
}

func TestOrderToMessage_encodings(t *testing.T) {
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
	orderType := abi.TypedData{
		Types: abi.Types{"TradeOrder": []abi.Type{
			{Name: "sender", Type: "address"},
			{Name: "subaccount", Type: "bytes32"},
			{Name: "quantity", Type: "uint128"},
			{Name: "price", Type: "uint128"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "side", Type: "uint8"},
			{Name: "engineType", Type: "uint8"},
			{Name: "productId", Type: "uint32"},
			{Name: "nonce", Type: "uint64"},
			{Name: "signedAt", Type: "uint64"},
		}},
	}
	message, err := order.ToMessage()
	if err != nil {
		t.Fatal(err)
	}
	check := func(typ, name string, messageKey string, wantHex string) {
		t.Helper()
		b, err := orderType.EncodePrimitiveValue(typ, message[messageKey], 2)
		if err != nil {
			t.Fatal(err)
		}
		if common.Bytes2Hex(b) != wantHex {
			t.Fatalf("%s: got %s want %s", name, common.Bytes2Hex(b), wantHex)
		}
	}
	check("address", "sender", "sender", "000000000000000000000000deadbeef00000000000000000000000000000000")
	check("bytes32", "subaccount", "subaccount", "123456789abcde00000000000000000000000000000000000000000000000000")
	check("uint128", "quantity", "quantity", "000000000000000000000000000000000000000000000000000000003b9aca00")
	check("uint128", "price", "price", "000000000000000000000000000000000000000000000000000002ba7def3000")
}
