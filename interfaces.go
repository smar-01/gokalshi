package gokalshi

import "context"

// HTTPClient defines the contract for the Kalshi HTTP API client.
// Consumers should type-annotate against this interface, not *Client directly.
type HTTPClient interface {
	// Account
	GetAccountAPILimits(ctx context.Context) (GetAccountApiLimitsResponse, error)
	GetAccountEndpointCosts(ctx context.Context) (GetAccountEndpointCostsResponse, error)
	UpgradeAPIUsageLevel(ctx context.Context) error

	// Exchange
	GetExchangeStatus(ctx context.Context) (ExchangeStatus, error)
	GetExchangeSchedule(ctx context.Context) (GetExchangeScheduleResponse, error)
	GetUserDataTimestamp(ctx context.Context) (GetUserDataTimestampResponse, error)
	GetSeriesFeeChanges(ctx context.Context, params GetSeriesFeeChangesParams) (GetSeriesFeeChangesResponse, error)

	// Orders
	GetOrder(ctx context.Context, orderID string) (GetOrderResponse, error)
	GetOrders(ctx context.Context, params GetOrdersParams) (GetOrdersResponse, error)
	GetQueuePositions(ctx context.Context, params GetQueuePositionsParams) (GetOrderQueuePositionsResponse, error)
	GetQueuePosition(ctx context.Context, orderID string) (GetOrderQueuePositionResponse, error)

	// Portfolio
	GetBalance(ctx context.Context) (GetBalanceResponse, error)
	GetBalanceWithParams(ctx context.Context, params GetBalanceParams) (GetBalanceResponse, error)
	GetPositions(ctx context.Context, params GetPositionsParams) (GetPositionsResponse, error)
	GetFills(ctx context.Context, params GetFillsParams) (GetFillsResponse, error)
	GetSettlements(ctx context.Context, params GetSettlementsParams) (GetSettlementsResponse, error)
	GetDeposits(ctx context.Context, params GetDepositsParams) (GetDepositsResponse, error)
	GetWithdrawals(ctx context.Context, params GetWithdrawalsParams) (GetWithdrawalsResponse, error)
	GetPortfolioRestingOrderTotalValue(ctx context.Context) (GetPortfolioRestingOrderTotalValueResponse, error)

	// Subaccounts
	CreateSubaccount(ctx context.Context) (CreateSubaccountResponse, error)
	GetSubaccountBalances(ctx context.Context) (GetSubaccountBalancesResponse, error)
	GetSubaccountNetting(ctx context.Context) (GetSubaccountNettingResponse, error)
	UpdateSubaccountNetting(ctx context.Context, req UpdateSubaccountNettingRequest) error
	ApplySubaccountTransfer(ctx context.Context, req ApplySubaccountTransferRequest) (ApplySubaccountTransferResponse, error)
	GetSubaccountTransfers(ctx context.Context, params GetSubaccountTransfersParams) (GetSubaccountTransfersResponse, error)

	// Order Groups
	CreateOrderGroup(ctx context.Context, req CreateOrderGroupRequest) (CreateOrderGroupResponse, error)
	GetOrderGroups(ctx context.Context, params GetOrderGroupsParams) (GetOrderGroupsResponse, error)
	GetOrderGroup(ctx context.Context, orderGroupID string, params GetOrderGroupParams) (GetOrderGroupResponse, error)
	DeleteOrderGroup(ctx context.Context, orderGroupID string, params DeleteOrderGroupParams) error
	ResetOrderGroup(ctx context.Context, orderGroupID string, params OrderGroupActionParams) error
	TriggerOrderGroup(ctx context.Context, orderGroupID string, params OrderGroupActionParams) error
	UpdateOrderGroupLimit(ctx context.Context, orderGroupID string, req UpdateOrderGroupLimitRequest, params UpdateOrderGroupLimitParams) error

	// Markets
	GetMarketOrderbook(ctx context.Context, ticker string, params GetOrderbookParams) (GetMarketOrderbookResponse, error)
	GetMarketOrderbooks(ctx context.Context, params GetMarketOrderbooksParams) (GetMarketOrderbooksResponse, error)
	GetTrades(ctx context.Context, params GetTradesParams) (GetTradesResponse, error)
	GetMarket(ctx context.Context, ticker string) (GetMarketResponse, error)
	GetMarkets(ctx context.Context, params GetMarketsParams) (GetMarketsResponse, error)
	GetMarketCandlesticks(ctx context.Context, seriesTicker, ticker string, params GetMarketCandlesticksParams) (GetMarketCandlesticksResponse, error)
	GetBatchMarketCandlesticks(ctx context.Context, params GetBatchMarketCandlesticksParams) (BatchGetMarketCandlesticksResponse, error)

	// Events
	GetEvent(ctx context.Context, eventTicker string, params GetEventParams) (GetEventResponse, error)
	GetEvents(ctx context.Context, params GetEventsParams) (GetEventsResponse, error)
	GetEventMetadata(ctx context.Context, eventTicker string) (GetEventMetadataResponse, error)
	GetMultivariateEvents(ctx context.Context, params GetMultivariateEventsParams) (GetMultivariateEventsResponse, error)
	GetEventCandlesticks(ctx context.Context, seriesTicker, eventTicker string, params GetEventCandlesticksParams) (GetEventCandlesticksResponse, error)
	GetEventForecastPercentileHistory(ctx context.Context, seriesTicker, eventTicker string, params GetEventForecastPercentileHistoryParams) (GetEventForecastPercentilesHistoryResponse, error)
	GetEventFeeChanges(ctx context.Context, params GetEventFeeChangesParams) (GetEventFeeChangesResponse, error)

	// Series
	GetSeries(ctx context.Context, seriesTicker string, params GetSeriesParams) (GetSeriesResponse, error)
	GetSeriesList(ctx context.Context, params GetSeriesListParams) (GetSeriesListResponse, error)

	// Search
	GetTagsByCategories(ctx context.Context) (GetTagsForSeriesCategoriesResponse, error)
	GetFiltersBySport(ctx context.Context) (GetFiltersBySportsResponse, error)

	// API Keys
	GetAPIKeys(ctx context.Context) (GetApiKeysResponse, error)
	CreateAPIKey(ctx context.Context, req CreateApiKeyRequest) (CreateApiKeyResponse, error)
	GenerateAPIKey(ctx context.Context, req GenerateApiKeyRequest) (GenerateApiKeyResponse, error)
	DeleteAPIKey(ctx context.Context, apiKey string) error

	// Event Orders (V2)
	CreateOrderV2(ctx context.Context, req CreateOrderV2Request) (CreateOrderV2Response, error)
	BatchCreateOrdersV2(ctx context.Context, req BatchCreateOrdersV2Request) (BatchCreateOrdersV2Response, error)
	BatchCancelOrdersV2(ctx context.Context, req BatchCancelOrdersV2Request) (BatchCancelOrdersV2Response, error)
	CancelOrderV2(ctx context.Context, orderID string, params CancelOrderV2Params) (CancelOrderV2Response, error)
	AmendOrderV2(ctx context.Context, orderID string, req AmendOrderV2Request) (AmendOrderV2Response, error)
	DecreaseOrderV2(ctx context.Context, orderID string, req DecreaseOrderV2Request) (DecreaseOrderV2Response, error)
	AmendOrderV2WithParams(ctx context.Context, orderID string, req AmendOrderV2Request, params OrderSubaccountParams) (AmendOrderV2Response, error)
	DecreaseOrderV2WithParams(ctx context.Context, orderID string, req DecreaseOrderV2Request, params OrderSubaccountParams) (DecreaseOrderV2Response, error)

	// Historical
	GetHistoricalCutoff(ctx context.Context) (GetHistoricalCutoffResponse, error)
	GetHistoricalFills(ctx context.Context, params GetHistoricalFillsParams) (GetFillsResponse, error)
	GetHistoricalOrders(ctx context.Context, params GetHistoricalOrdersParams) (GetOrdersResponse, error)
	GetHistoricalTrades(ctx context.Context, params GetHistoricalTradesParams) (GetTradesResponse, error)
	GetHistoricalMarkets(ctx context.Context, params GetHistoricalMarketsParams) (GetMarketsResponse, error)
	GetHistoricalMarket(ctx context.Context, ticker string) (GetMarketResponse, error)
	GetHistoricalMarketCandlesticks(ctx context.Context, ticker string, params GetHistoricalMarketCandlesticksParams) (GetMarketCandlesticksHistoricalResponse, error)

	// Incentive Programs
	GetIncentivePrograms(ctx context.Context) (GetIncentiveProgramsResponse, error)

	// Live Data
	GetLiveDataBatch(ctx context.Context, params GetLiveDataBatchParams) (GetLiveDatasResponse, error)
	GetLiveDataByMilestone(ctx context.Context, milestoneID string, params GetLiveDataParams) (GetLiveDataResponse, error)
	GetMilestoneGameStats(ctx context.Context, milestoneID string) (GetGameStatsResponse, error)
	GetLiveData(ctx context.Context, dataType, milestoneID string, params GetLiveDataParams) (GetLiveDataResponse, error)

	// Milestones
	GetMilestones(ctx context.Context, params GetMilestonesParams) (GetMilestonesResponse, error)
	GetMilestone(ctx context.Context, milestoneID string) (GetMilestoneResponse, error)

	// Multivariate Event Collections
	GetMultivariateEventCollections(ctx context.Context, params GetMultivariateEventCollectionsParams) (GetMultivariateEventCollectionsResponse, error)
	GetMultivariateEventCollection(ctx context.Context, collectionTicker string) (GetMultivariateEventCollectionResponse, error)
	CreateMarketInMultivariateEventCollection(ctx context.Context, collectionTicker string, req CreateMarketInMultivariateEventCollectionRequest) (CreateMarketInMultivariateEventCollectionResponse, error)

	// Structured Targets
	GetStructuredTargets(ctx context.Context, params GetStructuredTargetsParams) (GetStructuredTargetsResponse, error)
	GetStructuredTarget(ctx context.Context, structuredTargetID string) (GetStructuredTargetResponse, error)

	// Communications
	GetCommunicationsID(ctx context.Context) (GetCommunicationsIDResponse, error)
	CreateRFQ(ctx context.Context, req CreateRFQRequest) (CreateRFQResponse, error)
	GetRFQs(ctx context.Context, params GetRFQsParams) (GetRFQsResponse, error)
	GetRFQ(ctx context.Context, rfqID string) (GetRFQResponse, error)
	DeleteRFQ(ctx context.Context, rfqID string) error
	CreateQuote(ctx context.Context, req CreateQuoteRequest) (CreateQuoteResponse, error)
	GetQuotes(ctx context.Context, params GetQuotesParams) (GetQuotesResponse, error)
	GetQuote(ctx context.Context, quoteID string) (GetQuoteResponse, error)
	DeleteQuote(ctx context.Context, quoteID string) error
	AcceptQuote(ctx context.Context, quoteID string, req AcceptQuoteRequest) error
	ConfirmQuote(ctx context.Context, quoteID string) error

	// Lifecycle
	ConfigureRateLimits(ctx context.Context) error
	Close()
}

// WebSocketClient defines the contract for the Kalshi WebSocket client.
type WebSocketClient interface {
	Connect(ctx context.Context) error
	Connected() bool
	Close() error
	MsgCh() <-chan []byte
	ListenLoop(ctx context.Context)
	Subscribe(ctx context.Context, channels []string, tickers []string, opts ...SubscribeOption) error
	Unsubscribe(ctx context.Context, channels []string) error
	AddMarkets(ctx context.Context, tickers []string, channels []string, opts ...SubscribeOption) error
	RemoveMarkets(ctx context.Context, tickers []string, channels []string) error
	GetSnapshot(ctx context.Context, tickers []string, channels []string) error
}

// Compile-time interface satisfaction checks.
var (
	_ HTTPClient      = (*Client)(nil)
	_ WebSocketClient = (*WSClient)(nil)
)
