package restClient

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/qiwi1272/ethereal-go"
)

const USER_AGENT = "ethereal-go-sdk/1.1.0"

type Client struct {
	baseURL    string
	http       *http.Client
	Subaccount *ethereal.Subaccount
	Types      *abi.TypedData
	pk         *ecdsa.PrivateKey
	Address    string
}

type Environment string

const (
	Testnet Environment = "https://api.etherealtest.net"
	Mainnet Environment = "https://api.ethereal.trade"
)

func NewClient(ctx context.Context, pk string, env Environment) (*Client, error) {
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
		ExpectContinueTimeout: 0,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
	}

	client := &Client{
		baseURL: string(env),
		http: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}

	// load pk
	if pk == "" {
		return nil, errors.New("no private key provided; ETHEREAL_PK not set in environment")
	}

	// parse key, set address
	var err error
	client.pk, err = crypto.HexToECDSA(strip0x(pk))
	if err != nil {
		return nil, errors.New("unable to parse private key, likely invalid format")
	}
	client.Address = crypto.PubkeyToAddress(client.pk.PublicKey).Hex()
	// ethereal env setup
	_, err = client.InitDomain(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("unable to compute domain hash: "), err)
	}

	if err := client.InitSubaccount(ctx); err != nil {
		return nil, errors.Join(errors.New("failed to fetch subaccount: "), err)
	}

	return client, nil
}

// ---------- REST ----------

type Response[T any] struct {
	Data T `json:"data"`
}

type SignedGenericMessage struct {
	Data      any    `json:"data"`
	Signature string `json:"signature"`
}

func (e *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, e.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out := new(bytes.Buffer)
	_, err = out.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ethereal error %d: %s", resp.StatusCode, out.String())
	}
	return out.Bytes(), nil
}

// ---------- Setup ----------
func (e *Client) InitDomain(ctx context.Context) (string, error) {
	// init eip 712 data from rpc
	data, err := e.do(ctx, "GET", "/v1/rpc/config", nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		Domain   abi.TypedDataDomain `json:"domain"`
		SigTypes map[string]string   `json:"signatureTypes"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return "", err
	}

	// parse flattened type data
	parsedTypes := abi.Types{}
	for primaryType, schema := range resp.SigTypes {
		types, err := ParseTypeSchema(schema)
		if err != nil {
			return "", err
		}
		parsedTypes[primaryType] = types
	}
	// hardcode domain type
	parsedTypes["EIP712Domain"] = []abi.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	e.Types = &abi.TypedData{
		Types:  parsedTypes,
		Domain: resp.Domain,
	}

	domain, err := e.Types.HashStruct("EIP712Domain", e.Types.Domain.Map())
	if err != nil {
		panic("failed to compute domain hash: " + err.Error())
	}
	domainHash = domain
	return hex.EncodeToString(domain), nil
}

func (e *Client) InitSubaccount(ctx context.Context) error {
	path := fmt.Sprintf("/v1/subaccount?sender=%s", e.Address)
	data, err := e.do(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	var resp Response[[]ethereal.Subaccount]
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	if len(resp.Data) == 0 {
		return errors.New("no subaccounts found")
	}
	e.Subaccount = &resp.Data[0] // NOTE: currently only one subaccount per client is supported

	return nil
}

// ---------- Methods ----------

func (e *Client) GetPosition(ctx context.Context) ([]Position, error) {
	path := fmt.Sprintf("/v1/position?subaccountId=%s&open=%v", e.Subaccount.Id, true)
	data, err := e.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]Position]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetAccountBalance(ctx context.Context) ([]*AccountBalance, error) {
	path := fmt.Sprintf("/v1/subaccount/balance?subaccountId=%s", e.Subaccount.Id)
	data, err := e.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]*AccountBalance]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetProductMap(ctx context.Context) (map[string]Product, error) {
	data, err := e.do(ctx, "GET", "/v1/product", nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]Product]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	products := make(map[string]Product)

	for _, p := range resp.Data {
		products[p.Ticker] = p
	}

	return products, nil
}

func (e *Client) BatchOrder(ctx context.Context, orders []*Order) ([]OrderCreated, error) {
	payload := make([]Signable, len(orders))
	for i, order := range orders {
		payload[i] = order
	}
	return SendBatch[OrderCreated](ctx, e, Create, payload)
}
