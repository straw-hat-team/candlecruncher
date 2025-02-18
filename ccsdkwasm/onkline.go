package ccsdkwasm

import (
	"encoding/json"

	"github.com/extism/go-pdk"
	"github.com/straw-hat-team/candlecruncher/ccsdk"
)

type OnKline[State any] func(ccsdk.OnKlineInput[State]) (*ccsdk.OnKlineOutput[State], error)

func NewOnKline[State any](handler OnKline[State]) int32 {
	var input ccsdk.OnKlinePdkInput
	err := json.Unmarshal(pdk.Input(), &input)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	var state State
	err = json.Unmarshal(input.State, &state)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	output, err := handler(ccsdk.OnKlineInput[State]{
		Symbol:      input.Symbol,
		Kline:       input.Kline,
		StartedTime: input.StartedTime,
		State:       &state,
	})
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	newStateBytes, err := json.Marshal(output.State)
	if err != nil {
		pdk.SetError(err)
		return -1
	}
	outputBytes, err := json.Marshal(ccsdk.OnKlinePdkOutput{
		ActionType: output.ActionType,
		State:      newStateBytes,
	})
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.Output(outputBytes)
	return 0
}

func NewOnKlineFromStrategy[State any, Params any](strategy *ccsdk.NativeStrategy[State, Params]) int32 {
	return NewOnKline[State](OnKline[State](strategy.OnKlineFunc))
}
