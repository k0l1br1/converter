package main

import (
	"errors"
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

	u, d, err := bins.Convert(cStg, check, intChan, opts.Symbol)
	if err != nil {
		if errors.Is(err, bins.ErrInterrupted) {
			fmt.Println("Interrupted!")
			return exitInterrupt
		}
		errorPrint(err)
		return exitError
	}

	fmt.Printf("All Done. Up bins: %d, Down bins: %d\n", u, d)
	return exitOk
}

func main() {
	os.Exit(run())
}
