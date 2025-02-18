package domain_test

import (
	"testing"

	"github.com/straw-hat-team/candlecruncher/domain"

	"github.com/stretchr/testify/require"
)

func TestTimeframe1m(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730419260000)

	hasPassed, err := start.TimeframeHasPassed(end, domain.Timeframe1m)
	require.NoError(t, err)
	require.Equal(t, true, hasPassed)
}

func TestTimeframe3m(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730419320000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe5m(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730419440000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe5m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe15m(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730420040000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe30m(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730420940000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe1h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730422740000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe2h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730426340000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe4h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730433540000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe6h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730440740000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe8h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730447940000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe12h(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730462340000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe12h,
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe1d(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730505540000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe1d,
		domain.Timeframe12h,
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe3d(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1730678340000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe3d,
		domain.Timeframe1d,
		domain.Timeframe12h,
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe1w(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1731023940000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe1w,
		domain.Timeframe3d,
		domain.Timeframe1d,
		domain.Timeframe12h,
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}

func TestTimeframe1M(t *testing.T) {
	start := domain.OpenTime(1730419200000)
	end := domain.OpenTime(1733011140000)

	for _, tf := range []domain.Timeframe{
		domain.Timeframe1M,
		domain.Timeframe1w,
		domain.Timeframe3d,
		domain.Timeframe1d,
		domain.Timeframe12h,
		domain.Timeframe8h,
		domain.Timeframe6h,
		domain.Timeframe4h,
		domain.Timeframe2h,
		domain.Timeframe1h,
		domain.Timeframe30m,
		domain.Timeframe15m,
		domain.Timeframe5m,
		domain.Timeframe3m,
		domain.Timeframe1m,
	} {
		hasPassed, err := start.TimeframeHasPassed(end, tf)
		require.NoError(t, err)
		require.Equal(t, true, hasPassed)
	}
}
