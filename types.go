package cryptosignals

type SignalIndicator string
type Position string

type SignalEvent struct {
	Indicator SignalIndicator
	Position  Position
}

type Notificator func(event SignalEvent)

type Signal interface {
	Start() error
}

const (
	Long          Position        = "Long"
	Short         Position        = "Short"
	IndicatorMACD SignalIndicator = "MACD"
)
