package ccsdk

type OnInitialStateInput[Params any] struct {
	Params Params
}

type OnInitialStateOutput[State any] struct {
	State State
}

type OnInitialState[State any, Params any] func(OnInitialStateInput[Params]) (*OnInitialStateOutput[State], error)
