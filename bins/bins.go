package bins

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
	"path/filepath"
)

const (
	DefaultFilePerm = 0644
	DefaultDirPerm  = 0744
	DefaultDataDir  = "bins"
	DefaultBinsExt  = ".bin"
	DefaultCheckExt = ".json"
	BinByteSize     = 9
)

type Bin struct {
	IsUp   uint8
	Time   uint32
	Volume float32
}

func InitNew(path string) error {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, DefaultFilePerm)
	if err != nil {
		return err
	}
	fd.Close()
	return nil
}

func InitNewDefault(symbol string) error {
	path, err := defaultDataPath()
	if err != nil {
		return err
	}
	return InitNew(filepath.Join(path, symbol+DefaultBinsExt))
}

func Write(path string, b []Bin) error {
	if len(b) == 0 {
		return nil
	}
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, DefaultFilePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	bs := make([]byte, BinByteSize)
	for i := range b {
		bs[0] = b[i].IsUp
		binary.LittleEndian.PutUint32(bs[1:5], b[i].Time)
		binary.LittleEndian.PutUint32(bs[5:9], math.Float32bits(b[i].Volume))
		// write one bin
		if _, err := fd.Write(bs); err != nil {
			return err
		}
	}
	return nil
}

func WriteDefault(symbol string, b []Bin) error {
	path, err := defaultDataPath()
	if err != nil {
		return err
	}
	return Write(filepath.Join(path, symbol+DefaultBinsExt), b)
}

func ReadAll(path string) ([]Bin, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(bs)%BinByteSize != 0 {
		return nil, errors.New("corrupted bins file")
	}
	b := make([]Bin, int(len(bs)/BinByteSize))
	n := len(b)
	var off int
	for i := 0; i < n; i++ {
		off = i * BinByteSize
		b[i].IsUp = bs[off]
		b[i].Time = binary.LittleEndian.Uint32(bs[1+off : 5+off])
		b[i].Volume = math.Float32frombits(binary.LittleEndian.Uint32(bs[5+off : 9+off]))
	}
	return b, nil
}

func ReadAllDefault(symbol string) ([]Bin, error) {
	path, err := defaultDataPath()
	if err != nil {
		return nil, err
	}
	return ReadAll(filepath.Join(path, symbol+DefaultBinsExt))
}

func defaultDataPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(wd, DefaultDataDir)
	if err := os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return "", err
	}
	return dir, nil
}
