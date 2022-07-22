package cryptosignals

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type signalMACD struct {
	startSeries    *techan.TimeSeries
	lastTimeSeries int64
	symbol         string
	interval       string
	lastSignal     big.Decimal
	notified       bool
	indicator      SignalIndicator
	stop           chan struct{}
	_stop          chan struct{}
	res            chan SignalEvent
}

func (s *signalMACD) klineHandler(event *futures.WsKlineEvent) {
	if s.lastTimeSeries == 0 {
		s.lastTimeSeries = event.Kline.EndTime
	} else {
		if s.lastTimeSeries != event.Kline.EndTime {
			s.lastTimeSeries = event.Kline.EndTime
			time.Sleep(time.Second)
			s.init()
		}
	}
	if s.notified {
		return
	}
	ser := techan.NewTimeSeries()
	ser.Candles = s.startSeries.Candles
	s.parseKline(futures.Kline{
		OpenTime:  event.Kline.StartTime,
		Open:      event.Kline.Open,
		High:      event.Kline.High,
		Low:       event.Kline.Low,
		Volume:    event.Kline.Volume,
		CloseTime: event.Kline.EndTime,
		Close:     event.Kline.Close,
		TradeNum:  event.Kline.TradeNum,
	}, ser)
	signal := s.getSignal(ser)
	cross, way := s.hasCross(signal)
	if cross {
		s.res <- SignalEvent{Indicator: s.indicator, Position: way, Symbol: s.symbol}
		s.notified = true
	}
}

func (s *signalMACD) parseKline(kline futures.Kline, series *techan.TimeSeries) {
	start := time.Unix(0, 1000000*kline.OpenTime)
	end := time.Unix(0, 1000000*kline.CloseTime)
	period := techan.NewTimePeriod(start, end.Add(1*time.Millisecond).Sub(start))

	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(kline.Open)
	candle.ClosePrice = big.NewFromString(kline.Close)
	candle.MaxPrice = big.NewFromString(kline.High)
	candle.MinPrice = big.NewFromString(kline.Low)
	candle.Volume = big.NewFromString(kline.Volume)
	candle.TradeCount = uint(kline.TradeNum)
	series.AddCandle(candle)
}

func (s *signalMACD) getSeries() (*techan.TimeSeries, error) {
	binanceClient := futures.NewClient("", "")
	klines, err := binanceClient.
		NewKlinesService().
		Symbol(s.symbol).
		Interval(s.interval).
		Limit(500).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	series := techan.NewTimeSeries()
	serLen := len(klines)
	for index, kline := range klines {
		if index+1 == serLen {
			break
		}
		s.parseKline(*kline, series)
	}
	return series, nil
}

func (s *signalMACD) getSignal(series *techan.TimeSeries) big.Decimal {
	closePrices := techan.NewClosePriceIndicator(series)
	macd := techan.NewMACDIndicator(closePrices, 12, 26)
	hist := techan.NewMACDHistogramIndicator(macd, 9)
	signal := hist.Calculate(series.LastIndex())
	return signal
}

func (s *signalMACD) hasCross(signal big.Decimal) (bool, Position) {
	last := s.lastSignal.Float()
	curr := signal.Float()
	if curr < 0 && last > 0 {
		return true, Short
	}
	if last < 0 && curr > 0 {
		return true, Long
	}
	return false, ""
}

func (s *signalMACD) init() error {
	var err error
	s.startSeries, err = s.getSeries()
	if err != nil {
		return err
	}
	s.lastSignal = s.getSignal(s.startSeries)
	s.notified = false
	return err
}

func (s *signalMACD) errHandler(err error) {
	select {
	case s.stop <- struct{}{}:
		s._stop <- struct{}{}
		close(s._stop)
	default:
	}
}

// Start run indicator
// non block func
func (s *signalMACD) Start() (res chan SignalEvent, stop chan struct{}, err error) {
	err = s.init()
	s.res = make(chan SignalEvent)
	s.stop = make(chan struct{})
	res = s.res
	if err != nil {
		return res, stop, err
	}
	_, stop, err = futures.WsKlineServe(s.symbol, s.interval, s.klineHandler, s.errHandler)
	s._stop = stop
	return res, s.stop, err
}
