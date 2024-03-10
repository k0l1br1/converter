package main

import (
	"errors"
	"testing"
)

func TestOptionsParser(t *testing.T) {
	args := []string{"loader", "--symbol", "ethusdt"}
	_, err := parseOptions(args[1:])
	if err != nil && !errors.Is(err, errReqSymbol) {
		t.Errorf("parse options error %s", err.Error())
	}

	opts, _ := parseOptions(args)
	wantSymbol := "ETHUSDT"
	if opts.Symbol != wantSymbol {
		t.Errorf("parse symbol want %s, got %s", wantSymbol, opts.Symbol)
	}

	args = append(args, "--is-new")
	opts, _ = parseOptions(args)
	if !opts.IsNew {
		t.Error("parse is-new flag failing")
	}

	args = append(args, "--step", "0.1")
	opts, _ = parseOptions(args)
	var wantStep float32 = 0.1
	if opts.Step != wantStep {
		t.Errorf("parse step: want %.2f, got %.2f", wantStep, opts.Step)
	}
}
