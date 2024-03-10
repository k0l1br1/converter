package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const usage = `usage: converter -s <symbol> [options]
    -s, --symbol    The pair candles to convert to bins
    -n, --is-new    The flag to start new bins line for a symbol    
    --step          Percent step for bins (default 0.4%)
`

const DefaultStep = 0.4

var (
	errReqSymbol   = errors.New("symbol is required")
	errInvalidStep = errors.New("step must be more then 0.1")
)

func help() {
	os.Stdout.WriteString(usage)
	os.Exit(exitOk)
}

type options struct {
	IsNew  bool
	Symbol string
	Step   float32
}

func validateOptions(opts *options) error {
	if opts.Symbol == "" {
		return errReqSymbol
	}
	if opts.Step < 0.1 {
		return errInvalidStep
	}
	return nil
}

func parseOptions(args []string) (*options, error) {
	if len(args) < 2 {
		help()
	}
	opts := &options{Step: DefaultStep}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			help()
		case "-n", "--is-new":
			opts.IsNew = true
		case "-s", "--symbol":
			j := i + 1
			// during validation there will be a check for an empty string
			// it is not necessary to check it here
			if len(args) > j && !strings.HasPrefix(args[j], "-") {
				opts.Symbol = strings.ToUpper(args[j])
				i++
			}
		case "--step":
			j := i + 1
			if len(args) > j && !strings.HasPrefix(args[j], "-") {
				v, err := strconv.ParseFloat(args[j], 32)
				if err != nil {
					return nil, err
				}
				opts.Step = float32(v)
				i++
			}
		}
	}

	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return opts, nil
}
