package ccsdk

import (
	"encoding/json"

	"github.com/straw-hat-team/candlecruncher/domain"
)

// PDK types for serialization
type State json.RawMessage

type OnKlinePdkInput struct {
	Symbol      domain.Symbol   `json:"symbol"`
	Kline       domain.Kline    `json:"kline"`
	State       State           `json:"state"`
	StartedTime domain.OpenTime `json:"started_time"`
}

type OnKlinePdkOutput struct {
	ActionType domain.ActionType `json:"type"`
	State      State             `json:"state"`
}

// Core types
type OnKline[State any] func(OnKlineInput[State]) (*OnKlineOutput[State], error)

type OnKlineInput[State any] struct {
	State       *State
	Symbol      domain.Symbol
	Kline       domain.Kline
	StartedTime domain.OpenTime
}

type OnKlineOutput[State any] struct {
	ActionType domain.ActionType
	State      State
}

type OnKlineTimeframes[State any] struct {
	Timeframe domain.Timeframe
	OnKline   OnKline[State]
}

// Helper methods
func (o OnKlineInput[State]) TimeframeHasPassed(timeframe domain.Timeframe) (bool, error) {
	return o.StartedTime.TimeframeHasPassed(o.Kline.OpenTime, timeframe)
}

func (o OnKlineInput[State]) WhenTimeframeHasPassed(timeframe domain.Timeframe, callback OnKline[State]) (*OnKlineOutput[State], error) {
	passed, err := o.StartedTime.TimeframeHasPassed(o.Kline.OpenTime, timeframe)
	if err != nil {
		return nil, err
	}

	if passed {
		return callback(o)
	}

	return nil, nil
}

// OnFirstKlineTimeframeMatch will call the first callback that matches the timeframe of the kline.
// If no callback matches, it will return a HoldActionType.
func (o OnKlineInput[State]) OnFirstKlineTimeframeMatch(callbacks ...OnKlineTimeframes[State]) (*OnKlineOutput[State], error) {
	for _, cfg := range callbacks {
		if cfg.Timeframe == domain.Timeframe1m {
			return cfg.OnKline(o)
		}

		// Calculate minutes difference between start and current kline
		diffMinutes := o.Kline.DiffMinutesSinceStart(o.StartedTime)

		// Get the timeframe minutes
		tfMinutes, err := cfg.Timeframe.Minutes()
		if err != nil {
			return nil, err
		}

		// Check if this is an exact match for the timeframe
		if diffMinutes > 0 && diffMinutes%tfMinutes == 0 {
			return cfg.OnKline(o)
		}
	}

	return &OnKlineOutput[State]{
		ActionType: domain.HoldActionType,
		State:      *o.State,
	}, nil
}
