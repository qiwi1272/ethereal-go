package rest

import (
	"context"
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

type Environment string

const (
	Testnet Environment = "https://api.etherealtest.net"
	Mainnet Environment = "https://api.ethereal.trade"
)

type Client struct {
	ethereal.RestClient
}

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

	client := &Client{RestClient: ethereal.RestClient{
		BaseURL: string(env),
		Http: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}}

	// load pk
	if pk == "" {
		return nil, errors.New("no private key provided; ETHEREAL_PK not set in environment")
	}

	// parse key, set address
	if len(pk) > 1 && pk[:2] == "0x" {
		pk = pk[2:]
	}
	if ecdsa, err := crypto.HexToECDSA(pk); err == nil {
		client.SetPk(ecdsa)
	} else {
		return nil, err
	}
	// ethereal env setup
	var err error
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

// ---------- Setup ----------
func (e *Client) InitDomain(ctx context.Context) (string, error) {
	// init eip 712 data from rpc
	data, err := e.Do(ctx, "GET", "/v1/rpc/config", nil)
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
		types, err := ethereal.ParseTypeSchema(schema)
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

	types := &abi.TypedData{
		Types:  parsedTypes,
		Domain: resp.Domain,
	}

	e.SetTypes(types)

	domain, err := types.HashStruct("EIP712Domain", types.Domain.Map())
	if err != nil {
		panic("failed to compute domain hash: " + err.Error())
	}
	ethereal.DomainHash = domain
	return hex.EncodeToString(domain), nil
}

func (e *Client) InitSubaccount(ctx context.Context) error {
	path := fmt.Sprintf("/v1/subaccount?sender=%s", e.Address)
	data, err := e.Do(ctx, "GET", path, nil)
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

func (e *Client) GetPosition(ctx context.Context) ([]ethereal.Position, error) {
	path := fmt.Sprintf("/v1/position?subaccountId=%s&open=%v", e.Subaccount.Id, true)
	data, err := e.Do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]ethereal.Position]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetAccountBalance(ctx context.Context) ([]*ethereal.AccountBalance, error) {
	path := fmt.Sprintf("/v1/subaccount/balance?subaccountId=%s", e.Subaccount.Id)
	data, err := e.Do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]*ethereal.AccountBalance]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetProductMap(ctx context.Context) (map[string]ethereal.Product, error) {
	data, err := e.Do(ctx, "GET", "/v1/product", nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]ethereal.Product]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	products := make(map[string]ethereal.Product)

	for _, p := range resp.Data {
		products[p.Ticker] = p
	}

	return products, nil
}

func (e *Client) PlaceOrders(ctx context.Context, orders []*ethereal.Order) ([]*ethereal.OrderCreated, error) {
	batch := ethereal.NewBatch[*ethereal.Order, *ethereal.OrderCreated](orders...)
	var cl ethereal.BatchOrderClient = e
	return batch.SendBatch(ctx, &cl, ethereal.Create)
}
