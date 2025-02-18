package main

import (
	"github.com/straw-hat-team/candlecruncher/ccsdkwasm"
	"github.com/straw-hat-team/candlecruncher/strategies/fuxa/fuxastrategy"
)

//export onInitialState
func onInitialState() int32 {
	return ccsdkwasm.NewOnInitialState(fuxastrategy.OnInitialState)
}

//export onKline
func onKline() int32 {
	return ccsdkwasm.NewOnKline(fuxastrategy.OnKline)
}

func main() {}
