package backtestingruntime

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/straw-hat-team/candlecruncher/ccsdk"
	"github.com/straw-hat-team/candlecruncher/domain"

	extism "github.com/extism/go-sdk"
)

// RuntimeMetrics tracks the performance metrics
type RuntimeMetrics struct {
	TotalTrades              int         `json:"total_trades"`
	WinningTrades            int         `json:"winning_trades"`
	LosingTrades             int         `json:"losing_trades"`
	MaxConsecutiveWins       int         `json:"max_consecutive_wins"`
	MaxConsecutiveLosses     int         `json:"max_consecutive_losses"`
	currentConsecutiveWins   int         `json:"-"` // Internal tracking
	currentConsecutiveLosses int         `json:"-"` // Internal tracking
	TotalProfits             TotalProfit `json:"total_profits"`
	TotalLosses              TotalLoss   `json:"total_losses"`
	TotalBarsInTrades        int         `json:"total_bars_in_trades"`
	TotalBarsInWinning       int         `json:"total_bars_in_winning"`
	TotalBarsInLosing        int         `json:"total_bars_in_losing"`
}

type BacktestRuntimeState struct {
	LastProcessedTime  domain.CloseTime
	LastProcessedIndex uint64
	Position           *domain.Position
	State              ccsdk.State
	ExecutorParams     domain.ExecutorParams
}

func (s BacktestRuntimeState) ShouldClosePosition(auction domain.ActionType, kline domain.Kline) bool {
	if s.Position == nil {
		return false
	}
	return s.ExecutorParams.CloseMode.ClosePosition(*s.Position, auction, kline)
}

type BacktestRuntime struct {
	mu sync.RWMutex

	state       *BacktestRuntimeState
	metrics     RuntimeMetrics
	plugin      *extism.Plugin
	StartedTime *domain.OpenTime
}

func NewBacktestRuntime(
	ctx context.Context,
	wasm []extism.Wasm,
	params domain.ExecutorParams,
	initialStateParams interface{}) (*BacktestRuntime, error) {
	plugin, err := extism.NewPlugin(
		ctx,
		extism.Manifest{Wasm: wasm},
		extism.PluginConfig{EnableWasi: true},
		[]extism.HostFunction{})

	if err != nil {
		return nil, err
	}

	initialStateParamsBytes, err := json.Marshal(initialStateParams)
	if err != nil {
		return nil, err
	}

	_, initialState, err := plugin.Call("onInitialState", initialStateParamsBytes)
	if err != nil {
		return nil, err
	}

	return &BacktestRuntime{
		state: &BacktestRuntimeState{
			Position:           nil,
			State:              initialState,
			ExecutorParams:     params,
			LastProcessedIndex: 0,
		},
		metrics: RuntimeMetrics{},
		plugin:  plugin,
	}, nil
}

func (rt *BacktestRuntime) GetMetrics() RuntimeMetrics {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	return rt.metrics
}

func (rt *BacktestRuntime) handleStrategyCloseMode(symbol domain.Symbol, kline domain.Kline, actionType domain.ActionType) {
	if rt.state.ShouldClosePosition(actionType, kline) {
		rt.closePosition(kline)
		return
	}

	if actionType == domain.BuyActionType {
		if rt.state.Position == nil && rt.state.ExecutorParams.StartPosition != domain.ShortStartPosition {
			rt.state.Position = &domain.Position{
				Symbol:   symbol,
				Quantity: rt.state.ExecutorParams.Quantity,
				Kline:    kline,
				Kind:     domain.LongPosition,
			}
		}
	} else if actionType == domain.SellActionType {
		if rt.state.Position == nil && rt.state.ExecutorParams.StartPosition != domain.LongStartPosition {
			rt.state.Position = &domain.Position{
				Symbol:   symbol,
				Quantity: rt.state.ExecutorParams.Quantity,
				Kline:    kline,
				Kind:     domain.ShortPosition,
			}
		}
	}
}

func (rt *BacktestRuntime) handleTPSLCloseMode(symbol domain.Symbol, kline domain.Kline, actionType domain.ActionType) {
	if actionType == domain.BuyActionType {
		if rt.state.Position == nil && rt.state.ExecutorParams.StartPosition != domain.ShortStartPosition {
			rt.state.Position = &domain.Position{
				Symbol:   symbol,
				Quantity: rt.state.ExecutorParams.Quantity,
				Kline:    kline,
				Kind:     domain.LongPosition,
			}
			return
		}
	} else if actionType == domain.SellActionType {
		if rt.state.Position == nil && rt.state.ExecutorParams.StartPosition != domain.LongStartPosition {
			rt.state.Position = &domain.Position{
				Symbol:   symbol,
				Quantity: rt.state.ExecutorParams.Quantity,
				Kline:    kline,
				Kind:     domain.ShortPosition,
			}
			return
		}
	}

	if rt.state.ShouldClosePosition(actionType, kline) {
		rt.closePosition(kline)
	}
}

func (rt *BacktestRuntime) closePosition(kline domain.Kline) {
	profitLoss := kline.Close.Sub(rt.state.Position.Kline.Close).ToPrice()
	barsInTrade := int(kline.OpenTime - rt.state.Position.Kline.OpenTime)

	rt.metrics.TotalTrades++
	rt.metrics.TotalBarsInTrades += barsInTrade

	if profitLoss.GreaterThanZero() {
		rt.metrics.WinningTrades++
		rt.metrics.TotalProfits = rt.metrics.TotalProfits.AddPriceQuantity(profitLoss, rt.state.Position.Quantity)
		rt.metrics.TotalBarsInWinning += barsInTrade

		rt.metrics.currentConsecutiveWins++
		rt.metrics.currentConsecutiveLosses = 0
		if rt.metrics.currentConsecutiveWins > rt.metrics.MaxConsecutiveWins {
			rt.metrics.MaxConsecutiveWins = rt.metrics.currentConsecutiveWins
		}
	} else {
		rt.metrics.LosingTrades++
		rt.metrics.TotalLosses = rt.metrics.TotalLosses.RemovePriceQuantity(profitLoss, rt.state.Position.Quantity)
		rt.metrics.TotalBarsInLosing += barsInTrade

		rt.metrics.currentConsecutiveLosses++
		rt.metrics.currentConsecutiveWins = 0
		if rt.metrics.currentConsecutiveLosses > rt.metrics.MaxConsecutiveLosses {
			rt.metrics.MaxConsecutiveLosses = rt.metrics.currentConsecutiveLosses
		}
	}

	rt.state.Position = nil
}

func (rt *BacktestRuntime) ProcessKline(symbol domain.Symbol, kline domain.Kline) error {

	rt.mu.Lock()
	defer rt.mu.Unlock()

	if rt.StartedTime == nil {
		rt.StartedTime = &kline.OpenTime
	}

	inputBytes, err := json.Marshal(ccsdk.OnKlinePdkInput{
		Symbol:      symbol,
		Kline:       kline,
		State:       rt.state.State,
		StartedTime: *rt.StartedTime,
	})
	if err != nil {
		return err
	}

	_, onKlineOutputBytes, err := rt.plugin.Call("onKline", inputBytes)
	if err != nil {
		return err
	}

	var onKlineOutput ccsdk.OnKlinePdkOutput
	if err := json.Unmarshal(onKlineOutputBytes, &onKlineOutput); err != nil {
		return err
	}

	rt.state.LastProcessedTime = kline.CloseTime
	rt.state.State = onKlineOutput.State

	switch rt.state.ExecutorParams.CloseMode.Type() {
	case domain.StrategyCloseModeType:
		rt.handleStrategyCloseMode(symbol, kline, onKlineOutput.ActionType)
	case domain.TPSLCloseModeType:
		rt.handleTPSLCloseMode(symbol, kline, onKlineOutput.ActionType)
	}

	rt.state.LastProcessedIndex++
	return nil
}
