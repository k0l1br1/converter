package bins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/valyala/fastjson"
)

const checkTmpl = `{
    "step": %f,
    "next_up": %.9f,
    "next_down": %.9f,
    "last_time": %d
}
`

type Checkpoint struct {
	step     float32
	NextUp   float32
	NextDown float32
	LastTime uint32
}

func NewCheckpoint(price, step float32, time uint32) *Checkpoint {
	check := Checkpoint{
		step:     step,
		NextUp:   price + price*step/100,
		NextDown: price - price*step/100,
		LastTime: time,
	}
	return &check
}

func LoadCheckpoint(path string) (*Checkpoint, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("parse checkpoint: %w", err)
	}

	check := Checkpoint{
		step:     float32(v.GetFloat64("step")),
		NextUp:   float32(v.GetFloat64("next_up")),
		NextDown: float32(v.GetFloat64("next_down")),
		LastTime: uint32(v.GetUint("last_time")),
	}
	return &check, nil
}

func LoadDefaultCheckpoint(symbol string) (*Checkpoint, error) {
	path, err := defaultDataPath()
	if err != nil {
		return nil, err
	}
	return LoadCheckpoint(filepath.Join(path, symbol+DefaultCheckExt))
}

func (c *Checkpoint) StepUp(time uint32) {
	dl := c.NextUp * c.step / 100
	c.NextDown = c.NextUp - dl
	c.NextUp += dl
	c.LastTime = time
}

func (c *Checkpoint) StepDown(time uint32) {
	dl := c.NextDown * c.step / 100
	c.NextUp = c.NextDown + dl
	c.NextDown -= dl
	c.LastTime = time
}

func (c *Checkpoint) Save(path string) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultFilePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	fmt.Fprintf(fd, checkTmpl, c.step, c.NextUp, c.NextDown, c.LastTime)
	return nil
}

func (c *Checkpoint) SaveDefault(symbol string) error {
	path, err := defaultDataPath()
	if err != nil {
		return err
	}
	return c.Save(filepath.Join(path, symbol+DefaultCheckExt))
}
