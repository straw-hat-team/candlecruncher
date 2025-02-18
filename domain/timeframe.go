package domain

import "errors"

type Timeframe string

const (
	Timeframe1m  Timeframe = "1m"
	Timeframe3m  Timeframe = "3m"
	Timeframe5m  Timeframe = "5m"
	Timeframe15m Timeframe = "15m"
	Timeframe30m Timeframe = "30m"
	Timeframe1h  Timeframe = "1h"
	Timeframe2h  Timeframe = "2h"
	Timeframe4h  Timeframe = "4h"
	Timeframe6h  Timeframe = "6h"
	Timeframe8h  Timeframe = "8h"
	Timeframe12h Timeframe = "12h"
	Timeframe1d  Timeframe = "1d"
	Timeframe3d  Timeframe = "3d"
	Timeframe1w  Timeframe = "1w"
	Timeframe1M  Timeframe = "1M"
)

var ErrTimeframeMinutesUnregistered = errors.New("timeframe minutes unregistered")

// Minutes returns the number of minutes for a given timeframe
func (t Timeframe) Minutes() (int64, error) {
	switch t {
	case Timeframe1m:
		return 1, nil
	case Timeframe3m:
		return 3, nil
	case Timeframe5m:
		return 5, nil
	case Timeframe15m:
		return 15, nil
	case Timeframe30m:
		return 30, nil
	case Timeframe1h:
		return 60, nil
	case Timeframe2h:
		return 120, nil
	case Timeframe4h:
		return 240, nil
	case Timeframe6h:
		return 360, nil
	case Timeframe8h:
		return 480, nil
	case Timeframe12h:
		return 720, nil
	case Timeframe1d:
		return 1440, nil
	case Timeframe3d:
		return 4320, nil
	case Timeframe1w:
		return 10080, nil
	case Timeframe1M:
		return 43200, nil // 1 month = 30 days
	default:
		return 0, ErrTimeframeMinutesUnregistered
	}
}
