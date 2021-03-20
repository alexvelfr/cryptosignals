package cryptosignals

// NewSignal create new signal for symbol, interval and indicator type
// symbol must contain only letters, example BTCUSDT
// interval spec https://binance-docs.github.io/apidocs/spot/en/#kline-candlestick-streams
func NewSignal(symbol, interval string, indicator SignalIndicator) Signal {
	switch indicator {
	case IndicatorMACD:
		return &signalMACD{symbol: symbol, interval: interval, indicator: indicator}
	}
	return nil
}
