package backtestingruntime

import (
	"github.com/straw-hat-team/candlecruncher/domain"

	"github.com/shopspring/decimal"
)

type TotalProfit domain.Price

func NewTotalProfitFromPrice(p domain.Price) TotalProfit {
	return TotalProfit(p)
}

func (t TotalProfit) ToPrice() domain.Price {
	return domain.Price(t)
}

func (t TotalProfit) ToDecimal() decimal.Decimal {
	return t.ToPrice().ToDecimal()
}

func (t TotalProfit) InexactFloat64() float64 {
	return t.ToPrice().ToDecimal().InexactFloat64()
}

func (t TotalProfit) AddPriceQuantity(price domain.Price, quantity domain.Quantity) TotalProfit {
	return NewTotalProfitFromPrice(t.ToPrice().Add(price.MulQuantity(quantity)))
}
