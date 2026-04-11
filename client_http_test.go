package etherealRest

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed testdata/rpc_config.json
var testRPCConfigJSON []byte

//go:embed testdata/subaccount_response.json
var testSubaccountFixture []byte

func TestSubaccount_fixtureJSON(t *testing.T) {
	var resp Response[[]Subaccount]
	if err := json.Unmarshal(testSubaccountFixture, &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Id == "" {
		t.Fatalf("unexpected fixture: %+v", resp.Data)
	}
}

func TestClient_Do_success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/ok" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"hello":"world"}`))
	}))
	defer ts.Close()

	cl := &Client{
		BaseURL: ts.URL,
		Http:    ts.Client(),
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}
	body, err := cl.Do(context.Background(), http.MethodGet, "/ok", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"hello":"world"}` {
		t.Fatalf("body: %s", body)
	}
}

func TestClient_Do_errorStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	}))
	defer ts.Close()

	cl := &Client{
		BaseURL: ts.URL,
		Http:    ts.Client(),
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}
	_, err := cl.Do(context.Background(), http.MethodGet, "/x", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ethereal error 400") || !strings.Contains(err.Error(), "nope") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestClient_Do_contextCancel(t *testing.T) {
	block := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	defer close(block)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cl := &Client{
		BaseURL: ts.URL,
		Http: &http.Client{
			Timeout: 2 * time.Second,
		},
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}
	_, err := cl.Do(ctx, http.MethodGet, "/slow", nil)
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestInitDomain_setsTypesAndDomainHash(t *testing.T) {
	prev := DomainHash
	t.Cleanup(func() { DomainHash = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/rpc/config":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(testRPCConfigJSON)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	cl := &Client{
		BaseURL: ts.URL,
		Http:    ts.Client(),
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}

	domainHex, err := cl.InitDomain(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(DomainHash) != 32 {
		t.Fatalf("DomainHash len: %d", len(DomainHash))
	}
	if domainHex == "" {
		t.Fatal("empty domain hex")
	}
	// Golden: InitDomain over testdata/rpc_config.json (update if fixture changes).
	const goldenDomainHex = "67c7a53d8e16034d15f3e8db7d7f152926c03b8caecc0bb351e842de909bf02e"
	if domainHex != goldenDomainHex {
		t.Fatalf("domain hash drift: got %s (update goldenDomainHex if testdata/rpc_config.json changed)", domainHex)
	}

	if cl.GetTypes() == nil || cl.GetTypes().Types["TradeOrder"] == nil {
		t.Fatal("TradeOrder types missing")
	}
	if cl.GetTypes().Types["CancelOrder"] == nil {
		t.Fatal("CancelOrder types missing")
	}
}

func TestInitSubaccount_firstEntry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/subaccount") {
			sender := r.URL.Query().Get("sender")
			resp := Response[[]Subaccount]{
				Data: []Subaccount{
					{Id: "0x01", Name: "0x02", Account: sender},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cl := &Client{
		BaseURL: ts.URL,
		Http:    ts.Client(),
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}
	if err := cl.InitSubaccount(context.Background()); err != nil {
		t.Fatal(err)
	}
	if cl.GetSubaccount().Account != cl.account.Address {
		t.Fatalf("subaccount account: got %q want %q", cl.GetSubaccount().Account, cl.account.Address)
	}
}

func TestInitSubaccount_emptyErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer ts.Close()

	cl := &Client{
		BaseURL: ts.URL,
		Http:    ts.Client(),
		account: NewSigner(mustECDSA(t, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a")),
	}
	if err := cl.InitSubaccount(context.Background()); err == nil {
		t.Fatal("expected error for empty subaccounts")
	}
}

func TestNewClient_httptest(t *testing.T) {
	prev := DomainHash
	t.Cleanup(func() { DomainHash = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/rpc/config":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(testRPCConfigJSON)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/subaccount"):
			sender := r.URL.Query().Get("sender")
			resp := Response[[]Subaccount]{
				Data: []Subaccount{
					{Id: "0x1111111111111111111111111111111111111111111111111111111111111111",
						Name:    "0x2222222222222222222222222222222222222222222222222222222222222222",
						Account: sender},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	pk := "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a"
	cl, err := NewClient(context.Background(), pk, Environment(ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	if cl.GetSubaccount().Account != cl.account.Address {
		t.Fatal("subaccount not wired")
	}
}

func mustECDSA(t *testing.T, hexKey string) *ecdsa.PrivateKey {
	t.Helper()
	k, err := crypto.HexToECDSA(strings.TrimPrefix(hexKey, "0x"))
	if err != nil {
		t.Fatal(err)
	}
	return k
}
