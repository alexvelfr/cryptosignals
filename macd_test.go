package cryptosignals

import (
	"fmt"
	"testing"
	"time"
)

func TestSignalMACD(t *testing.T) {
	s := NewSignal("BTCUSDT", "1m", IndicatorMACD, notify)
	stop, err := s.Start()
	if err != nil {
		panic(err)
	}
	<-stop
}

func notify(event SignalEvent) {
	fmt.Printf("%+v %s\n", event, time.Now().String())
}
