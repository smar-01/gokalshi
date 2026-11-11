package gokalshi

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureServer answers every request with body and records the method, path and query of each.
func captureServer(t *testing.T, body string) (*httptest.Server, *[]*http.Request) {
	t.Helper()
	var seen []*http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Clone(context.Background()))
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv, &seen
}

func onlyQuery(t *testing.T, seen []*http.Request) url.Values {
	t.Helper()
	require.Len(t, seen, 1)
	return seen[0].URL.Query()
}

func TestGetBalanceWithParams_Subaccount(t *testing.T) {
	srv, seen := captureServer(t, `{"balance":100,"balance_dollars":"1.00","portfolio_value":0,"updated_ts":1}`)
	c := newTestClient(t, srv.URL)
	resp, err := c.GetBalanceWithParams(context.Background(), GetBalanceParams{Subaccount: 3})
	require.NoError(t, err)
	assert.Equal(t, int64(100), resp.Balance)
	q := onlyQuery(t, *seen)
	assert.Equal(t, "3", q.Get("subaccount"))
	assert.False(t, q.Has("exchange_index"))
	assert.Equal(t, "/trade-api/v2/portfolio/balance", (*seen)[0].URL.Path)

	srv, seen = captureServer(t, `{"balance":1,"portfolio_value":0,"updated_ts":1}`)
	c = newTestClient(t, srv.URL)
	idx := 0
	_, err = c.GetBalanceWithParams(context.Background(), GetBalanceParams{ExchangeIndex: &idx})
	require.NoError(t, err)
	q = onlyQuery(t, *seen)
	assert.False(t, q.Has("subaccount"), "the primary account is Kalshi's default and is omitted")
	assert.Equal(t, "0", q.Get("exchange_index"))
}

func TestAmendAndDecreaseWithParams_Subaccount(t *testing.T) {
	srv, seen := captureServer(t, `{"order":{"order_id":"o1"}}`)
	c := newTestClient(t, srv.URL)
	_, err := c.AmendOrderV2WithParams(context.Background(), "o1", AmendOrderV2Request{Ticker: "T", Side: "bid", Price: "0.50", Count: "10"}, OrderSubaccountParams{Subaccount: 5})
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, (*seen)[0].Method)
	assert.Equal(t, "/trade-api/v2/portfolio/events/orders/o1/amend", (*seen)[0].URL.Path)
	assert.Equal(t, "5", onlyQuery(t, *seen).Get("subaccount"))

	srv, seen = captureServer(t, `{"order_id":"o1","remaining_count":"5"}`)
	c = newTestClient(t, srv.URL)
	by := "5"
	_, err = c.DecreaseOrderV2WithParams(context.Background(), "o1", DecreaseOrderV2Request{ReduceBy: &by}, OrderSubaccountParams{Subaccount: 5})
	require.NoError(t, err)
	assert.Equal(t, "/trade-api/v2/portfolio/events/orders/o1/decrease", (*seen)[0].URL.Path)
	assert.Equal(t, "5", onlyQuery(t, *seen).Get("subaccount"))

	srv, seen = captureServer(t, `{"order_id":"o1","remaining_count":"5"}`)
	c = newTestClient(t, srv.URL)
	_, err = c.DecreaseOrderV2WithParams(context.Background(), "o1", DecreaseOrderV2Request{ReduceBy: &by}, OrderSubaccountParams{})
	require.NoError(t, err)
	assert.False(t, onlyQuery(t, *seen).Has("subaccount"), "primary: same request as DecreaseOrderV2")
}

func TestOnlySubaccount_SendsZero(t *testing.T) {
	zero := 0
	srv, seen := captureServer(t, `{"orders":[],"cursor":""}`)
	c := newTestClient(t, srv.URL)
	_, err := c.GetOrders(context.Background(), GetOrdersParams{OnlySubaccount: &zero})
	require.NoError(t, err)
	assert.Equal(t, "0", onlyQuery(t, *seen).Get("subaccount"), "primary only, not every sub-account")

	srv, seen = captureServer(t, `{"orders":[],"cursor":""}`)
	c = newTestClient(t, srv.URL)
	_, err = c.GetOrders(context.Background(), GetOrdersParams{})
	require.NoError(t, err)
	assert.False(t, onlyQuery(t, *seen).Has("subaccount"), "unchanged: omitted means every sub-account")

	srv, seen = captureServer(t, `{"fills":[],"cursor":""}`)
	c = newTestClient(t, srv.URL)
	_, err = c.GetFills(context.Background(), GetFillsParams{OnlySubaccount: &zero})
	require.NoError(t, err)
	assert.Equal(t, "0", onlyQuery(t, *seen).Get("subaccount"))

	two := 2
	srv, seen = captureServer(t, `{"settlements":[],"cursor":""}`)
	c = newTestClient(t, srv.URL)
	_, err = c.GetSettlements(context.Background(), GetSettlementsParams{OnlySubaccount: &two})
	require.NoError(t, err)
	assert.Equal(t, "2", onlyQuery(t, *seen).Get("subaccount"))
}

// A whole request signed with an Ed25519 key: the headers on the wire verify against the public key over
// timestamp + method + path (no query).
func TestClient_Ed25519SignsRealRequest(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	srv, seen := captureServer(t, `{"balance":1,"portfolio_value":0,"updated_ts":1}`)
	cfg := testClientConfig(t, srv.URL)
	cfg.Credentials = NewEd25519Credentials("ed-key", priv)
	c, err := NewClient(cfg)
	require.NoError(t, err)
	_, err = c.GetBalanceWithParams(context.Background(), GetBalanceParams{Subaccount: 1})
	require.NoError(t, err)
	r := (*seen)[0]
	assert.Equal(t, "ed-key", r.Header.Get("KALSHI-ACCESS-KEY"))
	sig, err := base64.StdEncoding.DecodeString(r.Header.Get("KALSHI-ACCESS-SIGNATURE"))
	require.NoError(t, err)
	msg := r.Header.Get("KALSHI-ACCESS-TIMESTAMP") + r.Method + r.URL.Path
	assert.True(t, ed25519.Verify(pub, []byte(msg), sig))
}
