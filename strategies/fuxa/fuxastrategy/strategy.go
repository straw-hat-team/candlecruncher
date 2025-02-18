package fuxastrategy

import (
	"github.com/straw-hat-team/candlecruncher/ccsdk"
	"github.com/straw-hat-team/candlecruncher/domain"
)

var Strategy = ccsdk.NewNativeStrategy(OnInitialState, OnKline)

type Params struct {
}

type State struct {
}

func OnInitialState(input ccsdk.OnInitialStateInput[Params]) (*ccsdk.OnInitialStateOutput[State], error) {
	return &ccsdk.OnInitialStateOutput[State]{
		State: State{},
	}, nil
}

func OnKline(input ccsdk.OnKlineInput[State]) (*ccsdk.OnKlineOutput[State], error) {
	return input.OnFirstKlineTimeframeMatch(
		ccsdk.OnKlineTimeframes[State]{Timeframe: domain.Timeframe5m, OnKline: On5minKline},
		ccsdk.OnKlineTimeframes[State]{Timeframe: domain.Timeframe1m, OnKline: On1minKline},
	)
}

func On1minKline(input ccsdk.OnKlineInput[State]) (*ccsdk.OnKlineOutput[State], error) {
	return &ccsdk.OnKlineOutput[State]{
		ActionType: domain.BuyActionType,
		State:      *input.State,
	}, nil
}

func On5minKline(input ccsdk.OnKlineInput[State]) (*ccsdk.OnKlineOutput[State], error) {
	return &ccsdk.OnKlineOutput[State]{
		ActionType: domain.SellActionType,
		State:      *input.State,
	}, nil
}
