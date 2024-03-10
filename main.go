package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/k0l1br1/converter/bins"
	"github.com/k0l1br1/loader/candles"
)

const (
	exitOk        = 0
	exitError     = 1
	exitInterrupt = 130
	candlesLimit  = 10000
)

func errorPrint(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
}

func errorWrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func run() int {
	opts, err := parseOptions(os.Args)
	if err != nil {
		errorPrint(err)
		switch err {
		case errReqSymbol:
			return exitOk
		default:
			return exitError
		}
	}

	cStg, err := candles.DefaultStorage(opts.Symbol)
	if err != nil {
		errorPrint(err)
		return exitError
	}
	defer cStg.Close()

	var check *bins.Checkpoint
	if opts.IsNew {
		tmp := make([]candles.Candle, 1)
		_, err = cStg.Read(tmp)
		if err != nil && err != io.EOF {
			errorPrint(err)
			return exitError
		}
		check = bins.NewCheckpoint(tmp[0].CPrice, opts.Step, tmp[0].CTime)
		bins.InitNewDefault(opts.Symbol)
	} else {
		check, err = bins.LoadDefaultCheckpoint(opts.Symbol)
		if err != nil {
			errorPrint(err)
			return exitError
		}
	}

	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt, syscall.SIGTERM)

	cs := make([]candles.Candle, candlesLimit)
	bs := make([]bins.Bin, 0, 1000)
	u := 0
	d := 0
	var volume float32
	for {
		n, err := cStg.Read(cs)

		if err != nil && err != io.EOF {
			errorPrint(err)
			return exitError
		}

		for i := 0; i < n; i++ {
			// skip candles previous processed
			if cs[i].CTime > check.LastTime {
				volume += cs[i].Volume
				if cs[i].HPrice > check.NextUp {
					bs = append(bs, bins.Bin{IsUp: 1, Time: cs[i].CTime - check.LastTime, Volume: volume})
					check.StepUp(cs[i].CTime)
					volume = 0
					u++
				} else if cs[i].LPrice < check.NextDown {
					bs = append(bs, bins.Bin{IsUp: 0, Time: cs[i].CTime - check.LastTime, Volume: volume})
					check.StepDown(cs[i].CTime)
					volume = 0
					d++
				}
			}
		}

		if n < len(cs) || err == io.EOF {
			if err = bins.WriteDefault(opts.Symbol, bs); err != nil {
				errorPrint(errorWrap("write bins", err))
				return exitError
			}
			if err = check.SaveDefault(opts.Symbol); err != nil {
				errorPrint(errorWrap("save checkpoint", err))
				return exitError
			}
			fmt.Printf("All Done. Up bins: %d, Down bins: %d\n", u, d)
			return exitOk
		}

		select {
		case <-intChan:
			fmt.Println("Interrupted!")
			return exitInterrupt
		default:
			// meaning that the selects never block
		}
	}
}

func main() {
	os.Exit(run())
}
