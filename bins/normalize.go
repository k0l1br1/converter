package bins

type NBin struct {
	IsUp   bool
	Time   float32
	Volume float32
}

func Normilize(bs []Bin, nbs []NBin) {
	if len(bs) != len(nbs) {
		panic("invalid bins lenght")
	}

	n := len(bs)
	tMax := bs[0].Time
	vMax := bs[0].Volume
	for i := 0; i < n; i++ {
		if bs[i].Time > tMax {
			tMax = bs[i].Time
		}
		if bs[i].Volume > vMax {
			vMax = bs[i].Volume
		}
	}

	tMaxFloat := float32(tMax)
	for i := 0; i < n; i++ {
		nbs[i].IsUp = bs[i].IsUp == 1
		nbs[i].Time = float32(bs[i].Time) / tMaxFloat
		nbs[i].Volume = bs[i].Volume / vMax
	}
}
