package ccsdkwasm

import (
	"encoding/json"

	"github.com/extism/go-pdk"
	"github.com/straw-hat-team/candlecruncher/ccsdk"
)

type OnInitialState[State any, Params any] func(ccsdk.OnInitialStateInput[Params]) (*ccsdk.OnInitialStateOutput[State], error)

func NewOnInitialState[State any, Params any](handler OnInitialState[State, Params]) int32 {
	var pdkInput Params
	err := json.Unmarshal(pdk.Input(), &pdkInput)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	state, err := handler(ccsdk.OnInitialStateInput[Params]{
		Params: pdkInput,
	})
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdkOutput, err := json.Marshal(state)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.Output(pdkOutput)
	return 0
}

func NewOnInitialStateFromStrategy[State any, Params any](strategy *ccsdk.NativeStrategy[State, Params]) int32 {
	return NewOnInitialState[State, Params](OnInitialState[State, Params](strategy.OnInitialStateFunc))
}
