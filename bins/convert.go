package bins

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/k0l1br1/loader/candles"
)

const (
	candlesLimit = 10000
)

var ErrInterrupted = errors.New("Interrupted")

func errorWrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func Convert(stg *candles.Storage, check *Checkpoint, intChan chan os.Signal, symbol string) (int, int, error) {
	cs := make([]candles.Candle, candlesLimit)
	bs := make([]Bin, 0, 1000)
	u := 0
	d := 0
	var volume float32
	for {
		n, err := stg.Read(cs)
		if err != nil && err != io.EOF {
			return u, d, err
		}

		for i := 0; i < n; i++ {
			// skip candles previous processed
			if cs[i].CTime > check.LastTime {
				volume += cs[i].Volume
				if cs[i].HPrice > check.NextUp {
					bs = append(bs, Bin{IsUp: 1, Time: cs[i].CTime - check.LastTime, Volume: volume})
					check.StepUp(cs[i].CTime)
					volume = 0
					u++
				} else if cs[i].LPrice < check.NextDown {
					bs = append(bs, Bin{IsUp: 0, Time: cs[i].CTime - check.LastTime, Volume: volume})
					check.StepDown(cs[i].CTime)
					volume = 0
					d++
				}
			}
		}

		if n < len(cs) || err == io.EOF {
			if err = WriteDefault(symbol, bs); err != nil {
				return u, d, errorWrap("write bins", err)
			}
			if err = check.SaveDefault(symbol); err != nil {
				return u, d, errorWrap("save checkpoint", err)
			}
			return u, d, nil
		}

		select {
		case <-intChan:
			return u, d, ErrInterrupted
		default:
			// meaning that the selects never block
		}
	}
}
