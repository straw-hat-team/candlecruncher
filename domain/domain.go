package domain

import (
	"github.com/shopspring/decimal"
)

type Quantity decimal.Decimal

func NewQuantityFromFloat(value float64) Quantity {
	return Quantity(decimal.NewFromFloat(value))
}

func (q Quantity) ToDecimal() decimal.Decimal {
	return decimal.Decimal(q)
}

type Symbol string

const (
	SymbolBTCUSD Symbol = "BTCUSD"
)

type Epoch uint64
type CloseTime Epoch

// OpenTime represents the time at which the candlestick opens
// - Uniquely identified by their open time.
type OpenTime Epoch

func (o OpenTime) Add(t OpenTime) int64 {
	return int64(o + t)
}

func (o OpenTime) Sub(t OpenTime) int64 {
	return int64(o - t)
}

func (o OpenTime) DiffMinutes(t OpenTime) int64 {
	return o.Sub(t) / 60000
}

func (o OpenTime) MinutesHasPassed(to OpenTime, minutes int64) bool {
	if minutes <= 0 {
		return false
	}

	// Convert minutes to milliseconds
	expectedMs := minutes * 60000
	actualMs := to.Sub(o)

	// Check if at least N minutes have passed
	return actualMs >= expectedMs
}

func (o OpenTime) TimeframeHasPassed(to OpenTime, timeframe Timeframe) (bool, error) {
	if timeframe == Timeframe1m {
		return true, nil
	}

	minutes, err := timeframe.Minutes()
	if err != nil {
		return false, err
	}
	return o.MinutesHasPassed(OpenTime(to), minutes-1), nil
}

type Open Price
type High Price
type Low Price
type Close Price

func (c Close) Sub(c2 Close) Close {
	return Close(c.ToPrice().Sub(c2.ToPrice()))
}

func (c Close) PercentageOf(p Close) decimal.Decimal {
	return c.ToDecimal().Div(p.ToDecimal()).Mul(decimal.NewFromInt(100))
}

func (c Close) ToPrice() Price {
	return Price(c)
}

func (c Close) ToDecimal() decimal.Decimal {
	return c.ToPrice().ToDecimal()
}

type NumberOfTrades uint64

type Volume float64
type BaseAssetVolume Volume
type TakerBuyVolume Volume
type TakerBuyBaseAssetVolume Volume

type Kline struct {
	OpenTime                OpenTime                `json:"open_time"`
	Open                    Open                    `json:"open"`
	High                    High                    `json:"high"`
	Low                     Low                     `json:"low"`
	Close                   Close                   `json:"close"`
	Volume                  Volume                  `json:"volume"`
	CloseTime               CloseTime               `json:"close_time"`
	BaseAssetVolume         BaseAssetVolume         `json:"base_asset_volume"`
	NumberOfTrades          NumberOfTrades          `json:"number_of_trades"`
	TakerBuyVolume          TakerBuyVolume          `json:"taker_buy_volume"`
	TakerBuyBaseAssetVolume TakerBuyBaseAssetVolume `json:"taker_buy_base_asset_volume"`
}

func (k Kline) DiffMinutesSinceStart(t OpenTime) int64 {
	return (k.OpenTime.Sub(t) / 60000) + 1
}

type ActionType string

const (
	BuyActionType  ActionType = "BUY"
	SellActionType ActionType = "SELL"
	HoldActionType ActionType = "HOLD"
)

type CloseModeType string

const (
	TPSLCloseModeType     CloseModeType = "TPSL"     // Take Profit/Stop Loss mode
	StrategyCloseModeType CloseModeType = "STRATEGY" // Let strategy decide when to close
)

type PositionKind string

const (
	LongPosition  PositionKind = "LONG"
	ShortPosition PositionKind = "SHORT"
)

type Position struct {
	Symbol   Symbol       `json:"symbol"`
	Quantity Quantity     `json:"quantity"`
	Kline    Kline        `json:"kline"`
	Kind     PositionKind `json:"kind"`
}

type QuantityMode string

const (
	CoinQuantityMode       QuantityMode = "COIN"       // Fixed amount in coin (e.g. 0.1 BTC)
	CurrencyQuantityMode   QuantityMode = "CURRENCY"   // Fixed amount in currency (e.g. 1000 USD)
	PercentageQuantityMode QuantityMode = "PERCENTAGE" // Percentage of available balance
)

type StartPosition string

const (
	LongStartPosition    StartPosition = "LONG"    // Only long positions
	ShortStartPosition   StartPosition = "SHORT"   // Only short positions
	NeutralStartPosition StartPosition = "NEUTRAL" // Both long and short positions
)

type TakeProfit decimal.Decimal

func NewTakeProfitFromFloat(value float64) TakeProfit {
	return TakeProfit(decimal.NewFromFloat(value))
}

func (t TakeProfit) ToDecimal() decimal.Decimal {
	return decimal.Decimal(t)
}

func (t TakeProfit) LessThanOrEqualDecimal(d decimal.Decimal) bool {
	return t.ToDecimal().LessThanOrEqual(d)
}

func NewTakeProfitFromDecimal(d decimal.Decimal) TakeProfit {
	return TakeProfit(d)
}

type StopLoss decimal.Decimal

func NewStopLossFromFloat(value float64) StopLoss {
	return StopLoss(decimal.NewFromFloat(value))
}

func (t StopLoss) ToDecimal() decimal.Decimal {
	return decimal.Decimal(t)
}

func (t StopLoss) GreaterThanOrEqualDecimal(d decimal.Decimal) bool {
	return t.ToDecimal().GreaterThanOrEqual(d)
}

func NewStopLossFromDecimal(d decimal.Decimal) StopLoss {
	return StopLoss(d)
}

type ExecutorParams struct {
	QuantityMode QuantityMode `json:"quantity_mode"`
	Quantity     Quantity     `json:"quantity"`
	//InitialBalance Price         `json:"initial_balance"`
	CloseMode     CloseMode     `json:"close_mode"`
	StartPosition StartPosition `json:"start_position"`
}

type CloseMode interface {
	Type() CloseModeType
	ClosePosition(position Position, actionType ActionType, kline Kline) bool
}

type TpSlCloseMode struct {
	TakeProfit TakeProfit `json:"take_profit"`
	StopLoss   StopLoss   `json:"stop_loss"`
}

func (c TpSlCloseMode) Type() CloseModeType {
	return TPSLCloseModeType
}

func (c TpSlCloseMode) ClosePosition(position Position, _ ActionType, kline Kline) bool {
	var percentage decimal.Decimal
	if position.Kind == LongPosition {
		percentage = kline.Close.PercentageOf(position.Kline.Close).Sub(decimal.NewFromInt(100))
	} else {
		percentage = position.Kline.Close.PercentageOf(kline.Close)
	}

	tp := c.TakeProfit.ToDecimal()

	// Check if we've hit take profit (positive percentage >= take profit level)
	if percentage.GreaterThanOrEqual(decimal.NewFromInt(0)) && percentage.GreaterThanOrEqual(tp) {
		return true
	}

	// Check if we've hit stop loss (negative percentage <= stop loss level)
	if percentage.LessThan(decimal.NewFromInt(0)) && percentage.Mul(decimal.NewFromInt(-1)).GreaterThanOrEqual(c.StopLoss.ToDecimal()) {
		return true
	}

	return false
}

type StrategyCloseMode struct {
	TakeProfit TakeProfit `json:"take_profit"`
	StopLoss   StopLoss   `json:"stop_loss"`
}

func (c StrategyCloseMode) Type() CloseModeType {
	return StrategyCloseModeType
}

func (c StrategyCloseMode) ClosePosition(position Position, actionType ActionType, _ Kline) bool {
	return (position.Kind == ShortPosition && actionType == BuyActionType) ||
		(position.Kind == LongPosition && actionType == SellActionType)
}
