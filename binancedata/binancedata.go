package binancedata

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/straw-hat-team/candlecruncher/domain"

	"github.com/shopspring/decimal"
)

const fieldsPerRecord = 12

type HistoricalData struct {
	Symbol    domain.Symbol
	Timeframe domain.Timeframe
	Klines    []domain.Kline
}

func ParseCSV(filePath string, symbol domain.Symbol, interval domain.Timeframe) (*HistoricalData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = fieldsPerRecord

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	klines := make([]domain.Kline, len(records)-1)

	for i, record := range records[1:] {
		openTime, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		open, err := decimal.NewFromString(record[1])
		if err != nil {
			return nil, err
		}
		high, err := decimal.NewFromString(record[2])
		if err != nil {
			return nil, err
		}
		low, err := decimal.NewFromString(record[3])
		if err != nil {
			return nil, err
		}
		closePrice, err := decimal.NewFromString(record[4])
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		closeTime, err := strconv.ParseUint(record[6], 10, 64)
		if err != nil {
			return nil, err
		}
		baseAssetVolume, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, err
		}
		numberOfTrades, err := strconv.ParseUint(record[8], 10, 64)
		if err != nil {
			return nil, err
		}
		takerBuyVolume, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			return nil, err
		}
		takerBuyBaseAssetVolume, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			return nil, err
		}

		klines[i] = domain.Kline{
			OpenTime:                domain.OpenTime(openTime),
			Open:                    domain.Open(open),
			High:                    domain.High(high),
			Low:                     domain.Low(low),
			Close:                   domain.Close(closePrice),
			Volume:                  domain.Volume(volume),
			CloseTime:               domain.CloseTime(closeTime),
			BaseAssetVolume:         domain.BaseAssetVolume(baseAssetVolume),
			NumberOfTrades:          domain.NumberOfTrades(numberOfTrades),
			TakerBuyVolume:          domain.TakerBuyVolume(takerBuyVolume),
			TakerBuyBaseAssetVolume: domain.TakerBuyBaseAssetVolume(takerBuyBaseAssetVolume),
		}
	}

	return &HistoricalData{
		Symbol:    symbol,
		Timeframe: interval,
		Klines:    klines,
	}, nil
}
