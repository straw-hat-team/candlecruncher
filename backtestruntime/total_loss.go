package backtestingruntime

import (
	"github.com/shopspring/decimal"
	"github.com/straw-hat-team/candlecruncher/domain"
)

type TotalLoss domain.Price

func NewTotalLossFromPrice(p domain.Price) TotalLoss {
	return TotalLoss(p)
}

func (t TotalLoss) ToPrice() domain.Price {
	return domain.Price(t)
}

func (t TotalLoss) ToDecimal() decimal.Decimal {
	return t.ToPrice().ToDecimal()
}

func (t TotalLoss) InexactFloat64() float64 {
	return t.ToPrice().ToDecimal().InexactFloat64()
}

func (t TotalLoss) RemovePriceQuantity(price domain.Price, quantity domain.Quantity) TotalLoss {
	return NewTotalLossFromPrice(t.ToPrice().Add(price.MulQuantity(quantity)).Abs())
}
