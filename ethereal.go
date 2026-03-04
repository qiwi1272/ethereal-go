package ethereal // import "ethereal-dev"

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const USER_AGENT = "ethereal-go/1.0.0dev"

type EtherealClient interface {
	RestClient
}

type RestClient struct {
	BaseURL    string
	Http       *http.Client
	Subaccount *Subaccount
	types      *abi.TypedData
	pk         *ecdsa.PrivateKey
	Address    string
}

func (r *RestClient) SetPk(p *ecdsa.PrivateKey) {
	r.pk = p
	r.Address = crypto.PubkeyToAddress(p.PublicKey).Hex()
	*p = ecdsa.PrivateKey{}
}

func (r *RestClient) getPk() *ecdsa.PrivateKey {
	return r.pk
}

func (r *RestClient) SetTypes(t *abi.TypedData) {
	r.types = t
}

func (r *RestClient) GetTypes() *abi.TypedData {
	return r.types
}

func (r *RestClient) GetSubaccount() *Subaccount {
	return r.Subaccount
}

func (e *RestClient) Do(ctx context.Context, method, path string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, e.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Http.Do(req)
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
