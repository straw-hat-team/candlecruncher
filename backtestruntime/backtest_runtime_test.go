package backtestingruntime_test

import (
	"context"
	"testing"

	"github.com/straw-hat-team/candlecruncher/backtestingruntime"
	"github.com/straw-hat-team/candlecruncher/binancedata"
	"github.com/straw-hat-team/candlecruncher/domain"
	"github.com/straw-hat-team/candlecruncher/strategies/fuxa/fuxastrategy"

	extism "github.com/extism/go-sdk"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
	data, err := binancedata.ParseCSV(
		"../binancedata/testdata/BTCUSDT-1m-2024-11.csv",
		domain.SymbolBTCUSD,
		domain.Timeframe1m,
	)
	require.NoError(t, err)
	require.Equal(t, 43200, len(data.Klines))
	wasm := []extism.Wasm{
		extism.WasmFile{
			Path: "../tmp/fuxa.wasm",
			Name: "fuxa",
		},
	}

	backtestRuntime, err := backtestingruntime.NewBacktestRuntime(
		context.Background(),
		wasm,
		domain.ExecutorParams{
			Quantity:      domain.NewQuantityFromFloat(0.1),
			QuantityMode:  domain.CoinQuantityMode,
			CloseMode:     domain.StrategyCloseMode{},
			StartPosition: domain.NeutralStartPosition,
		},
		fuxastrategy.Params{},
	)
	require.NoError(t, err)

	for _, kline := range data.Klines {
		err := backtestRuntime.ProcessKline(
			data.Symbol,
			kline,
		)
		require.NoError(t, err)
	}
}

func TestTP(t *testing.T) {
	data, err := binancedata.ParseCSV(
		"../fakedata/FCKUSDT-PROFIT-2024-12.csv",
		domain.SymbolBTCUSD,
		domain.Timeframe1m,
	)
	require.NoError(t, err)
	require.Equal(t, 19, len(data.Klines))
	wasm := []extism.Wasm{
		extism.WasmFile{
			Path: "../tmp/fuxa.wasm",
			Name: "fuxa",
		},
	}

	backtestRuntime, err := backtestingruntime.NewBacktestRuntime(
		context.Background(),
		wasm,
		domain.ExecutorParams{
			Quantity:     domain.NewQuantityFromFloat(1),
			QuantityMode: domain.CoinQuantityMode,
			CloseMode: domain.TpSlCloseMode{
				TakeProfit: domain.TakeProfit(decimal.NewFromFloat(15.0)),
				StopLoss:   domain.StopLoss(decimal.NewFromFloat(15.0)),
			},
			StartPosition: domain.NeutralStartPosition,
		},
		fuxastrategy.Params{},
	)
	require.NoError(t, err)

	for _, kline := range data.Klines {
		err := backtestRuntime.ProcessKline(
			data.Symbol,
			kline,
		)
		require.NoError(t, err)
	}

	metrics := backtestRuntime.GetMetrics()
	require.Equal(t, metrics.TotalTrades, 1)
	require.Equal(t, metrics.WinningTrades, 1)
	require.Equal(t, metrics.LosingTrades, 0)
	require.Equal(t, metrics.TotalProfits.InexactFloat64(), 10500.00)
	require.Equal(t, metrics.TotalLosses.InexactFloat64(), 0.00)
	require.Equal(t, metrics.MaxConsecutiveWins, 1)
	require.Equal(t, metrics.MaxConsecutiveLosses, 0)
}

func TestSL(t *testing.T) {
	data, err := binancedata.ParseCSV(
		"../fakedata/FCKUSDT-LOSS-2024-12.csv",
		domain.SymbolBTCUSD,
		domain.Timeframe1m,
	)
	require.NoError(t, err)
	require.Equal(t, 19, len(data.Klines))
	wasm := []extism.Wasm{
		extism.WasmFile{
			Path: "../tmp/fuxa.wasm",
			Name: "fuxa",
		},
	}

	backtestRuntime, err := backtestingruntime.NewBacktestRuntime(
		context.Background(),
		wasm,
		domain.ExecutorParams{
			Quantity:     domain.NewQuantityFromFloat(1),
			QuantityMode: domain.CoinQuantityMode,
			CloseMode: domain.TpSlCloseMode{
				TakeProfit: domain.NewTakeProfitFromFloat(15.0),
				StopLoss:   domain.NewStopLossFromFloat(15.0),
			},
			StartPosition: domain.NeutralStartPosition,
		},
		fuxastrategy.Params{},
	)
	require.NoError(t, err)

	for _, kline := range data.Klines {
		err := backtestRuntime.ProcessKline(
			data.Symbol,
			kline,
		)
		require.NoError(t, err)
	}

	metrics := backtestRuntime.GetMetrics()
	require.Equal(t, metrics.TotalTrades, 1)
	require.Equal(t, metrics.WinningTrades, 0)
	require.Equal(t, metrics.LosingTrades, 1)
	require.Equal(t, metrics.TotalProfits.InexactFloat64(), 0.00)
	require.Equal(t, metrics.TotalLosses.InexactFloat64(), 10500.00)
	require.Equal(t, metrics.MaxConsecutiveWins, 0)
	require.Equal(t, metrics.MaxConsecutiveLosses, 1)
}
