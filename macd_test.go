package cryptosignals

import (
	"fmt"
	"testing"
	"time"
)

func TestSignalMACD(t *testing.T) {
	s := NewSignal("BTCUSDT", "1m", IndicatorMACD)
	res, stop, err := s.Start()
	if err != nil {
		panic(err)
	}
	go func() {
		for r := range res {
			fmt.Printf("%+v %s\n", r, time.Now().String())
		}
	}()
	<-stop
}
