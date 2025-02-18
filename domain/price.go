package domain

import "github.com/shopspring/decimal"

type Price decimal.Decimal

func NewPriceFromFloat(value float64) Price {
	return Price(decimal.NewFromFloat(value))
}

func NewPriceFromDecimal(d decimal.Decimal) Price {
	return Price(d)
}

func (p Price) ToDecimal() decimal.Decimal {
	return decimal.Decimal(p)
}

func (p Price) Sub(c2 Price) Price {
	return Price(p.ToDecimal().Sub(c2.ToDecimal()))
}

func (p Price) GreaterThan(c2 Price) bool {
	return p.ToDecimal().GreaterThan(c2.ToDecimal())
}

func (p Price) GreaterThanZero() bool {
	return p.ToDecimal().GreaterThan(decimal.NewFromInt(0))
}

func (p Price) MulQuantity(q Quantity) Price {
	return NewPriceFromDecimal(p.ToDecimal().Mul(q.ToDecimal()))
}

func (p Price) Abs() Price {
	return NewPriceFromDecimal(p.ToDecimal().Abs())
}

func (p Price) Add(p2 Price) Price {
	return NewPriceFromDecimal(p.ToDecimal().Add(p2.ToDecimal()))
}
