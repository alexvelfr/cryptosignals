package cryptosignals

type SignalIndicator string
type Position string

type SignalEvent struct {
	Indicator SignalIndicator
	Position  Position
	Symbol    string
}

type Notificator func(event SignalEvent)

type Signal interface {
	Start() (stop chan struct{}, err error)
}

const (
	Long          Position        = "Long"
	Short         Position        = "Short"
	IndicatorMACD SignalIndicator = "MACD"
)
