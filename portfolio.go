package gokalshi

import (
	"context"
)

// GetBalance — Get Balance
//
// GET /trade-api/v2/portfolio/balance
//
// Returns the balance and portfolio value for a member. Both values include
// all exchange indexes unless `exchange_index` is provided. Pass `subaccount`
// to use a subaccount instead of the primary account. This endpoint also
// accepts API keys with the 'read::portfolio_balance' scope.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-balance
func (c *Client) GetBalance(ctx context.Context) (GetBalanceResponse, error) {
	return getJSON[GetBalanceResponse](c, ctx, pathPortfolio+"/balance", nil)
}

// GetBalanceParams are GetBalance's query parameters. Subaccount 0, the primary account, is Kalshi's default and is
// omitted; ExchangeIndex nil reads every exchange index.
type GetBalanceParams struct {
	Subaccount    int
	ExchangeIndex *int
}

func (p GetBalanceParams) toMap() map[string]string {
	return NewQuery().
		Int("subaccount", p.Subaccount).
		IntPtr("exchange_index", p.ExchangeIndex).
		Build()
}

// GetBalanceWithParams is GetBalance for one sub-account and, optionally, one exchange index.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-balance
func (c *Client) GetBalanceWithParams(ctx context.Context, params GetBalanceParams) (GetBalanceResponse, error) {
	return getJSON[GetBalanceResponse](c, ctx, pathPortfolio+"/balance", params.toMap())
}

// GetPositions — Get Positions
//
// GET /trade-api/v2/portfolio/positions
//
// Restricts the positions to those with any of following fields with non-zero
// values, as a comma separated list. The following values are accepted:
// position, total_traded
//
// See https://docs.kalshi.com/api-reference/portfolio/get-positions
func (c *Client) GetPositions(ctx context.Context, params GetPositionsParams) (GetPositionsResponse, error) {
	return getJSON[GetPositionsResponse](c, ctx, pathPortfolio+"/positions", params.toMap())
}

// GetFills — Get Fills
//
// GET /trade-api/v2/portfolio/fills
//
// Endpoint for getting all fills for the member. A fill is when a trade you
// have is matched. Fills that occurred before the historical cutoff are only
// available via `GET /historical/fills`. See [Historical
// Data](https://docs.kalshi.com/getting_started/historical_data) for details.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-fills
func (c *Client) GetFills(ctx context.Context, params GetFillsParams) (GetFillsResponse, error) {
	return getJSON[GetFillsResponse](c, ctx, pathPortfolio+"/fills", params.toMap())
}

// GetSettlements — Get Settlements
//
// GET /trade-api/v2/portfolio/settlements
//
// Endpoint for getting the member's settlements historical track.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-settlements
func (c *Client) GetSettlements(ctx context.Context, params GetSettlementsParams) (GetSettlementsResponse, error) {
	return getJSON[GetSettlementsResponse](c, ctx, pathPortfolio+"/settlements", params.toMap())
}

// GetDeposits — Get Deposits
//
// GET /trade-api/v2/portfolio/deposits
//
// Endpoint for getting the member's deposit history.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-deposits
func (c *Client) GetDeposits(ctx context.Context, params GetDepositsParams) (GetDepositsResponse, error) {
	return getJSON[GetDepositsResponse](c, ctx, pathPortfolio+"/deposits", params.toMap())
}

// GetWithdrawals — Get Withdrawals
//
// GET /trade-api/v2/portfolio/withdrawals
//
// Endpoint for getting the member's withdrawal history.
//
// See https://docs.kalshi.com/api-reference/portfolio/get-withdrawals
func (c *Client) GetWithdrawals(ctx context.Context, params GetWithdrawalsParams) (GetWithdrawalsResponse, error) {
	return getJSON[GetWithdrawalsResponse](c, ctx, pathPortfolio+"/withdrawals", params.toMap())
}

// ---------------------------------------------------------------------------
// Query parameter types
// ---------------------------------------------------------------------------

// GetPositionsParams holds optional query parameters for GetPositions.
type GetPositionsParams struct {
	Ticker      string
	EventTicker string
	CountFilter string
	Limit       int
	Cursor      string
	Subaccount  int
}

func (p GetPositionsParams) toMap() map[string]string {
	return NewQuery().
		String("ticker", p.Ticker).
		String("event_ticker", p.EventTicker).
		String("count_filter", p.CountFilter).
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Int("subaccount", p.Subaccount).
		Build()
}

// GetFillsParams holds optional query parameters for GetFills.
type GetFillsParams struct {
	Ticker     string
	OrderID    string
	MinTs      int64
	MaxTs      int64
	Limit      int
	Cursor     string
	Subaccount int
	// OnlySubaccount, when set, limits the read to that one sub-account, the primary (0) included. Subaccount
	// omits 0 from the query, and Kalshi reads an omitted subaccount as every sub-account.
	OnlySubaccount *int
}

func (p GetFillsParams) toMap() map[string]string {
	return NewQuery().
		String("ticker", p.Ticker).
		String("order_id", p.OrderID).
		Int64("min_ts", p.MinTs).
		Int64("max_ts", p.MaxTs).
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Int("subaccount", p.Subaccount).
		IntPtr("subaccount", p.OnlySubaccount).
		Build()
}

// GetSettlementsParams holds optional query parameters for GetSettlements.
type GetSettlementsParams struct {
	Ticker      string
	EventTicker string
	MinTs       int64
	MaxTs       int64
	Limit       int
	Cursor      string
	Subaccount  int
	// OnlySubaccount, when set, limits the read to that one sub-account, the primary (0) included. Subaccount
	// omits 0 from the query, and Kalshi reads an omitted subaccount as every sub-account.
	OnlySubaccount *int
}

func (p GetSettlementsParams) toMap() map[string]string {
	return NewQuery().
		String("ticker", p.Ticker).
		String("event_ticker", p.EventTicker).
		Int64("min_ts", p.MinTs).
		Int64("max_ts", p.MaxTs).
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Int("subaccount", p.Subaccount).
		IntPtr("subaccount", p.OnlySubaccount).
		Build()
}

// GetDepositsParams holds optional query parameters for GetDeposits.
type GetDepositsParams struct {
	Limit  int
	Cursor string
}

func (p GetDepositsParams) toMap() map[string]string {
	return NewQuery().
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Build()
}

// GetWithdrawalsParams holds optional query parameters for GetWithdrawals.
type GetWithdrawalsParams struct {
	Limit  int
	Cursor string
}

func (p GetWithdrawalsParams) toMap() map[string]string {
	return NewQuery().
		Int("limit", p.Limit).
		String("cursor", p.Cursor).
		Build()
}
