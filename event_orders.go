package gokalshi

import (
	"context"
	"fmt"
)

// Token costs for order operations, used when the endpoint-costs table is
// unavailable. Both are billed per item on the batch endpoints.
//
// These are the documented values, confirmed against the live table: an order
// create is the default 10, a cancel is 2. A successful ConfigureRateLimits
// overrides them, so they only apply on the degraded path — which is exactly
// when getting them right matters, since that path also falls back to
// Basic-tier bucket sizes.
const (
	costCreateOrder = 10.0
	costCancelOrder = 2.0
)

// CreateOrderV2 — Create Order (V2)
//
// POST /trade-api/v2/portfolio/events/orders
//
// Endpoint for submitting event-market orders using the V2 request/response
// shape (single-book `bid`/`ask` side and fixed-point dollar prices). The
// legacy `/portfolio/orders` endpoint will be deprecated no earlier than May
// 6, 2026 — clients should migrate to this path.
//
// See https://docs.kalshi.com/api-reference/orders/create-order-v2
func (c *Client) CreateOrderV2(ctx context.Context, req CreateOrderV2Request) (CreateOrderV2Response, error) {
	return postJSON[CreateOrderV2Response](c, ctx, pathEventOrders, req, 10.0)
}

// BatchCreateOrdersV2 — Batch Create Orders (V2)
//
// POST /trade-api/v2/portfolio/events/orders/batched
//
// Endpoint for submitting a batch of event-market orders using the V2
// request/response shape. The maximum batch size scales with your tier's write
// budget — see [Rate Limits and Tiers](/getting_started/rate_limits).
//
// See https://docs.kalshi.com/api-reference/orders/batch-create-orders-v2
func (c *Client) BatchCreateOrdersV2(ctx context.Context, req BatchCreateOrdersV2Request) (BatchCreateOrdersV2Response, error) {
	return doJSON[BatchCreateOrdersV2Response](c, ctx, "POST", pathEventOrders+"/batched",
		0, costCreateOrder, len(req.Orders), req, nil)
}

// BatchCancelOrdersV2 — Batch Cancel Orders (V2)
//
// DELETE /trade-api/v2/portfolio/events/orders/batched
//
// Endpoint for cancelling a batch of event-market orders using the V2 response
// shape. To auto-route a cancellation, provide its `market_ticker` and omit
// `exchange_index` or set it to `-1`. The maximum batch size scales with your
// tier's write budget — see [Rate Limits and
// Tiers](/getting_started/rate_limits).
//
// See https://docs.kalshi.com/api-reference/orders/batch-cancel-orders-v2
func (c *Client) BatchCancelOrdersV2(ctx context.Context, req BatchCancelOrdersV2Request) (BatchCancelOrdersV2Response, error) {
	// One logical call remains one HTTP request; oversized costs fail in Acquire.
	return doJSON[BatchCancelOrdersV2Response](c, ctx, "DELETE", pathEventOrders+"/batched",
		0, costCancelOrder, len(req.Orders), req, nil)
}

// CancelOrderV2 — Cancel Order (V2)
//
// DELETE /trade-api/v2/portfolio/events/orders/{order_id}
//
// Endpoint for cancelling event-market orders using the V2 response shape. To
// auto-route the cancellation, provide `market_ticker` and omit
// `exchange_index` or set it to `-1`. Returns `{order_id, client_order_id,
// reduced_by}` rather than a full order object.
//
// See https://docs.kalshi.com/api-reference/orders/cancel-order-v2
func (c *Client) CancelOrderV2(ctx context.Context, orderID string, params CancelOrderV2Params) (CancelOrderV2Response, error) {
	path := fmt.Sprintf("%s/%s", pathEventOrders, orderID)
	return doJSON[CancelOrderV2Response](c, ctx, "DELETE", path, 0, costCancelOrder, 1, nil, params.toMap())
}

// AmendOrderV2 — Amend Order (V2)
//
// POST /trade-api/v2/portfolio/events/orders/{order_id}/amend
//
// Endpoint for amending the price and/or max fillable count of an existing
// event-market order using the V2 request/response shape. The request `count`
// is the updated total/max fillable count, equal to already filled count plus
// desired resting remaining count. This behavior matches the v1 amend
// endpoints; only the request/response shape differs.
//
// See https://docs.kalshi.com/api-reference/orders/amend-order-v2
func (c *Client) AmendOrderV2(ctx context.Context, orderID string, req AmendOrderV2Request) (AmendOrderV2Response, error) {
	path := fmt.Sprintf("%s/%s/amend", pathEventOrders, orderID)
	return postJSON[AmendOrderV2Response](c, ctx, path, req, 10.0)
}

// DecreaseOrderV2 — Decrease Order (V2)
//
// POST /trade-api/v2/portfolio/events/orders/{order_id}/decrease
//
// Endpoint for decreasing the remaining count of an existing event-market
// order using the V2 request/response shape. Exactly one of `reduce_by` or
// `reduce_to` must be provided.
//
// See https://docs.kalshi.com/api-reference/orders/decrease-order-v2
func (c *Client) DecreaseOrderV2(ctx context.Context, orderID string, req DecreaseOrderV2Request) (DecreaseOrderV2Response, error) {
	path := fmt.Sprintf("%s/%s/decrease", pathEventOrders, orderID)
	return postJSON[DecreaseOrderV2Response](c, ctx, path, req, 10.0)
}

// OrderSubaccountParams selects the sub-account an amend or a decrease acts on (query parameter subaccount). 0, the
// primary account, is Kalshi's default and is omitted.
type OrderSubaccountParams struct {
	Subaccount int
}

func (p OrderSubaccountParams) toMap() map[string]string {
	return NewQuery().Int("subaccount", p.Subaccount).Build()
}

// AmendOrderV2WithParams is AmendOrderV2 for an order in a sub-account.
//
// See https://docs.kalshi.com/api-reference/orders/amend-order-v2
func (c *Client) AmendOrderV2WithParams(ctx context.Context, orderID string, req AmendOrderV2Request, params OrderSubaccountParams) (AmendOrderV2Response, error) {
	path := fmt.Sprintf("%s/%s/amend", pathEventOrders, orderID)
	return doJSON[AmendOrderV2Response](c, ctx, "POST", path, 0, 10.0, 1, req, params.toMap())
}

// DecreaseOrderV2WithParams is DecreaseOrderV2 for an order in a sub-account.
//
// See https://docs.kalshi.com/api-reference/orders/decrease-order-v2
func (c *Client) DecreaseOrderV2WithParams(ctx context.Context, orderID string, req DecreaseOrderV2Request, params OrderSubaccountParams) (DecreaseOrderV2Response, error) {
	path := fmt.Sprintf("%s/%s/decrease", pathEventOrders, orderID)
	return doJSON[DecreaseOrderV2Response](c, ctx, "POST", path, 0, 10.0, 1, req, params.toMap())
}

// CancelOrderV2Params are query parameters for CancelOrderV2.
type CancelOrderV2Params struct {
	// MarketTicker enables auto-routing when ExchangeIndex is omitted or -1.
	MarketTicker string
	Subaccount   int
	// ExchangeIndex selects the exchange instance. 0 is the event-contract
	// instance and -1 asks Kalshi to auto-route by market ticker, so it is a
	// pointer: both are meaningful values that a zero-valued int cannot express.
	ExchangeIndex *int
}

func (p CancelOrderV2Params) toMap() map[string]string {
	return NewQuery().
		String("market_ticker", p.MarketTicker).
		Int("subaccount", p.Subaccount).
		IntPtr("exchange_index", p.ExchangeIndex).
		Build()
}

// CancelAllOrders — Cancel All Orders
//
// DELETE /trade-api/v2/portfolio/events/orders
//
// Cancels all resting event-market orders for the authenticated Direct member
// across every exchange shard. If `subaccount` is omitted, matching orders may
// come from any subaccount. If it is provided, only orders for that subaccount
// are eligible. Newly placed orders may also be cancelled during the minute
// after the request.
//
// See https://docs.kalshi.com/api-reference/orders/cancel-all-orders
func (c *Client) CancelAllOrders(ctx context.Context, params CancelAllOrdersParams) error {
	_, err := c.do(ctx, "DELETE", pathEventOrders, 0, costCancelOrder, 1, nil, params.toMap())
	return err
}

// CancelAllOrdersParams are the query parameters for CancelAllOrders.
type CancelAllOrdersParams struct {
	// Subaccount scopes the cancel. 0 is the primary subaccount, 1-63 the
	// others. It is a pointer because omitting it is not the same as naming
	// subaccount 0: an omitted subaccount cancels resting orders from *every*
	// subaccount, across every exchange shard. A plain int could not express
	// "the primary subaccount only", and would silently widen the blast radius
	// to the whole account.
	Subaccount *int
}

func (p CancelAllOrdersParams) toMap() map[string]string {
	return NewQuery().
		IntPtr("subaccount", p.Subaccount).
		Build()
}
