package bins

import (
	"testing"
)

const (
	testCheckFile = "/tmp/test-checkpoint.json"
	step          = 10
	price         = 100
	time          = 0
)

func TestCheckpointAlgo(t *testing.T) {
	check := NewCheckpoint(price, step, time)
	var wantNextUp float32 = 110
	var wantNextDown float32 = 90
	if check.NextUp != wantNextUp {
		t.Errorf("init nextUp: want %f, got %f", wantNextUp, check.NextUp)
	}
	if check.NextDown != wantNextDown {
		t.Errorf("init nextDown: want %f, got %f", wantNextDown, check.NextDown)
	}

	reset(check)
	check.StepUp(1)
	if check.NextUp != wantNextUp {
		t.Errorf("stepUp nextUp: want %f, got %f", wantNextUp, check.NextUp)
	}
	if check.NextDown != wantNextDown {
		t.Errorf("stepUp nextDown: want %f, got %f", wantNextDown, check.NextDown)
	}

	reset(check)
	check.StepDown(2)
	if check.NextUp != wantNextUp {
		t.Errorf("stepDown nextUp: want %f, got %f", wantNextUp, check.NextUp)
	}
	if check.NextDown != wantNextDown {
		t.Errorf("stepDown nextDown: want %f, got %f", wantNextDown, check.NextDown)
	}
}

func TestCheckpoint(t *testing.T) {
	check := NewCheckpoint(price, step, time)
	check.LastTime = 1

	err := check.Save(testCheckFile)
	if err != nil {
		t.Errorf("save checkpoint: %s", err.Error())
	}

	check2, err := LoadCheckpoint(testCheckFile)
	if err != nil {
		t.Errorf("load checkpoint: %s", err.Error())
	}
	if *check != *check2 {
		t.Errorf("loaded checkpoint not equal: want %v, got %v", check, check2)
	}
}

func reset(check *Checkpoint) {
	check.NextUp = price
	check.NextDown = price
	check.LastTime = 1
}
