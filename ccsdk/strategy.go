package ccsdk

type Strategy interface {
	OnInitialState2(params any) ([]byte, error)
	OnKline2(input OnKlinePdkInput) (*OnKlinePdkOutput, error)
}

type NativeStrategy[State any, Params any] struct {
	OnInitialStateFunc OnInitialState[State, Params]
	OnKlineFunc        OnKline[State]
}

func (n NativeStrategy[State, Params]) OnInitialState2(params any) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (n NativeStrategy[State, Params]) OnKline2(input OnKlinePdkInput) (*OnKlinePdkOutput, error) {
	//TODO implement me
	panic("implement me")
}

func NewNativeStrategy[State any, Params any](
	onInitialState OnInitialState[State, Params],
	onKline OnKline[State],
) *NativeStrategy[State, Params] {
	return &NativeStrategy[State, Params]{
		OnInitialStateFunc: onInitialState,
		OnKlineFunc:        onKline,
	}
}
