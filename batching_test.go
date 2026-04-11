package etherealRest

import (
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed testdata/rpc_config.json
var batchTestRPCConfig []byte

func TestSendBatch_twoCreateOrders(t *testing.T) {
	prev := DomainHash
	t.Cleanup(func() { DomainHash = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/rpc/config":
			_, _ = w.Write(batchTestRPCConfig)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/subaccount"):
			sender := r.URL.Query().Get("sender")
			resp := Response[[]Subaccount]{
				Data: []Subaccount{{
					Id:      "0x1111111111111111111111111111111111111111111111111111111111111111",
					Name:    "0x2222222222222222222222222222222222222222222222222222222222222222",
					Account: sender,
				}},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/order":
			b, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(b), `"signature"`) {
				http.Error(w, "no signature", 400)
				return
			}
			_, _ = w.Write([]byte(`{"id":"oid1","clientOrderId":"","filled":"0","result":"ok"}`))
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

	orders := []*Order{
		NewRawOrder(ORDER_LIMIT, 1, PERPETUAL, 1, 100, false, BUY, TIF_GTD),
		NewRawOrder(ORDER_LIMIT, 1, PERPETUAL, 1, 101, false, BUY, TIF_GTD),
	}
	batch := NewOrderBatch(orders)
	out, err := batch.SendBatch(context.Background(), cl, Create, cl.account)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 || out[0].Id != "oid1" || out[1].Id != "oid1" {
		t.Fatalf("responses: %+v", out)
	}
}

func TestSendBatch_cancelIntent_usesCancelOrderPrimaryType(t *testing.T) {
	prev := DomainHash
	t.Cleanup(func() { DomainHash = prev })

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/rpc/config":
			_, _ = w.Write(batchTestRPCConfig)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/subaccount"):
			sender := r.URL.Query().Get("sender")
			resp := Response[[]Subaccount]{
				Data: []Subaccount{{
					Id:      "0x1111111111111111111111111111111111111111111111111111111111111111",
					Name:    "0x2222222222222222222222222222222222222222222222222222222222222222",
					Account: sender,
				}},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/order/cancel":
			_, _ = w.Write([]byte(`{"id":"cx","clientOrderId":"","result":"ok"}`))
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

	c1 := NewCancel("order-a")
	c1.Build(cl)

	if _, err := Sign(c1, "TradeOrder", cl.account); err == nil {
		t.Fatal("TradeOrder primary type must not apply to OrderCancel EIP-712 message")
	}
	if _, err := Sign(c1, "CancelOrder", cl.account); err != nil {
		t.Fatal(err)
	}

	batch := NewBatch[*OrderCancelled]([]Signable{c1})
	_, err = batch.SendBatch(context.Background(), cl, Cancel, cl.account)
	if err != nil {
		t.Fatal(err)
	}
}
