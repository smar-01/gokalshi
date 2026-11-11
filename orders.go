package gokalshi

import (
	"context"
	"fmt"
)

// GetOrder — Get Order
//
// GET /trade-api/v2/portfolio/orders/{order_id}
//
// Endpoint for getting a single order.
//
// See https://docs.kalshi.com/api-reference/orders/get-order
func (c *Client) GetOrder(ctx context.Context, orderID string) (GetOrderResponse, error) {
	path := fmt.Sprintf("%s/%s", pathOrders, orderID)
	return getJSON[GetOrderResponse](c, ctx, path, nil)
}

// GetOrders — Get Orders
//
// GET /trade-api/v2/portfolio/orders
//
// Restricts the response to orders that have a certain status: resting,
// canceled, or executed. Orders that have been canceled or fully executed
// before the historical cutoff are only available via `GET
// /historical/orders`. Resting orders will always be available through this
// endpoint. See [Historical
// Data](https://docs.kalshi.com/getting_started/historical_data) for details.
//
// See https://docs.kalshi.com/api-reference/orders/get-orders
func (c *Client) GetOrders(ctx context.Context, params GetOrdersParams) (GetOrdersResponse, error) {
	return getJSON[GetOrdersResponse](c, ctx, pathOrders, params.toMap())
}

// GetQueuePositions — Get Queue Positions for Orders
//
// GET /trade-api/v2/portfolio/orders/queue_positions
//
// Endpoint for getting queue positions for all resting orders. Queue position
// represents the number of contracts that need to be matched before an order
// receives a partial or full match, determined using price-time priority.
//
// See https://docs.kalshi.com/api-reference/orders/get-order-queue-positions
func (c *Client) GetQueuePositions(ctx context.Context, params GetQueuePositionsParams) (GetOrderQueuePositionsResponse, error) {
	return getJSON[GetOrderQueuePositionsResponse](c, ctx, pathOrders+"/queue_positions", params.toMap())
}

// GetQueuePosition — Get Order Queue Position
//
// GET /trade-api/v2/portfolio/orders/{order_id}/queue_position
//
// Endpoint for getting an order's queue position in the order book. This
// represents the amount of orders that need to be matched before this order
// receives a partial or full match. Queue position is determined using a
// price-time priority.
//
// See https://docs.kalshi.com/api-reference/orders/get-order-queue-position
func (c *Client) GetQueuePosition(ctx context.Context, orderID string) (GetOrderQueuePositionResponse, error) {
	path := fmt.Sprintf("%s/%s/queue_position", pathOrders, orderID)
	return getJSON[GetOrderQueuePositionResponse](c, ctx, path, nil)
}

// ---------------------------------------------------------------------------
// Query parameter types
// ---------------------------------------------------------------------------

// GetOrdersParams holds optional query parameters for GetOrders.
type GetOrdersParams struct {
	Ticker      string
	EventTicker string
	Status      OrderStatus
	Limit       int
	Cursor      string
	MinTs       int64
	MaxTs       int64
	Subaccount  int
	// OnlySubaccount, when set, limits the read to that one sub-account, the primary (0) included. Subaccount
	// omits 0 from the query, and Kalshi reads an omitted subaccount as every sub-account.
	OnlySubaccount *int
}

func (p GetOrdersParams) toMap() map[string]string {
	return NewQuery().
		String("ticker", p.Ticker).
		String("event_ticker", p.EventTicker).
		String("status", string(p.Status)).
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Int64("min_ts", p.MinTs).
		Int64("max_ts", p.MaxTs).
		Int("subaccount", p.Subaccount).
		IntPtr("subaccount", p.OnlySubaccount).
		Build()
}

// GetQueuePositionsParams holds optional query parameters for GetQueuePositions.
type GetQueuePositionsParams struct {
	MarketTickers string
	EventTicker   string
	Subaccount    int
}

func (p GetQueuePositionsParams) toMap() map[string]string {
	return NewQuery().
		String("market_tickers", p.MarketTickers).
		String("event_ticker", p.EventTicker).
		Int("subaccount", p.Subaccount).
		Build()
}
