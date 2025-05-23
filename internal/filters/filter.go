package filters

import "math"

type Resampler interface {
	Kernel(value float64) float64
}

func sinc(value float64) float64 {
	if value == 0 {
		return 1
	}
	return math.Sin(math.Pi*value) / (math.Pi * value)
}

type Lanczos struct {
	Radius int
}

func NewLanczos(radius int) *Lanczos {
	return &Lanczos{
		Radius: radius,
	}
}

func (l *Lanczos) Kernel(value float64) float64 {
	value = math.Abs(value)
	if value < float64(l.Radius) {
		return sinc(value) * sinc(value/float64(l.Radius))
	}
	return 0
}
